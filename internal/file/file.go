package file

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/VoyakinH/lokle_backend/config"
	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ctx_utils"
	"github.com/VoyakinH/lokle_backend/internal/pkg/hasher"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ioutils"
	"github.com/VoyakinH/lokle_backend/internal/pkg/middleware"
	"github.com/VoyakinH/lokle_backend/internal/user/usecase"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type FileManager struct {
	rootPath    string
	userUseCase usecase.IUserUsecase
	logger      logrus.Logger
}

func SetFileRouting(router *mux.Router, uu usecase.IUserUsecase, auth middleware.AuthMiddleware, logger logrus.Logger) {
	fileManager := FileManager{
		rootPath:    config.File.RootPath,
		userUseCase: uu,
		logger:      logger,
	}

	fileAPI := router.PathPrefix("/api/v1/file/").Subrouter()
	fileAPI.Handle("/upload", auth.WithAuth(http.HandlerFunc(fileManager.Upload))).Methods(http.MethodPost)
}

const MaxUploadSize = 5 << 20 // 5MB

func isImage(fileType string) bool {
	imgTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if imgTypes[fileType] {
		return true
	}
	return false
}

func (fm *FileManager) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx_utils.GetUser(ctx)
	if user == nil {
		fm.logger.Errorf("%s failed get ctx user with [status=%d]", r.URL, http.StatusForbidden)
		ioutils.SendError(w, http.StatusForbidden, "no auth")
		return
	}

	// validate max file size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		fm.logger.Errorf("%s file is too big [status=%d]", r.URL, http.StatusBadRequest)
		ioutils.SendError(w, http.StatusBadRequest, "file is too big")
		return
	}

	// find file
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		fm.logger.Errorf("%s can't to find file in req [status=%d]", r.URL, http.StatusBadRequest)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}
	defer file.Close()

	// validate file type
	buf := make([]byte, fileHeader.Size)
	file.Read(buf)
	fileType := http.DetectContentType(buf)
	if !isImage(fileType) {
		fm.logger.Errorf("%s forbidden file type %s [status=%d]", r.URL, fileType, http.StatusBadRequest)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	// find user dir path
	var userDirPath string
	var userRoleID uint64
	switch user.Role {
	case models.ParentRole:
		parent, status, err := fm.userUseCase.GetParentByID(ctx, user.ID)
		if err != nil || status != http.StatusOK {
			fm.logger.Errorf("%s failed get parent user [status=%d] [error=%s]", r.URL, status, err)
			ioutils.SendError(w, status, "internal")
			return
		}
		userDirPath = parent.DirPath
		userRoleID = parent.ID
	case models.ChildRole:
		child, status, err := fm.userUseCase.GetChildByID(ctx, user.ID)
		if err != nil || status != http.StatusOK {
			fm.logger.Errorf("%s failed get child user [status=%d] [error=%s]", r.URL, status, err)
			ioutils.SendError(w, status, "internal")
			return
		}
		userDirPath = child.DirPath
		userRoleID = child.ID
	default:
		fm.logger.Errorf("%s unknown role while getting dir path %s [status=%d]", r.URL, user.Role.String(), http.StatusInternalServerError)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}

	// create user dir path if one not exists
	if userDirPath == "" {
		hashedPathName, err := hasher.HashAndSalt(user.Email)
		if err != nil {
			fm.logger.Errorf("%s failed to create hash for user path with [status=%d] [error=%s]", r.URL, http.StatusInternalServerError, err)
			ioutils.SendError(w, http.StatusInternalServerError, "internal")
			return
		}
		// replace all '/' else we break path
		hashedPathName = strings.ReplaceAll(hashedPathName, "/", "")
		hashedPathName += "/"
		err = os.MkdirAll(fm.rootPath+hashedPathName, os.ModePerm)
		if err != nil {
			fm.logger.Errorf("%s failed to create user's dir with [status=%d] [error=%s]", r.URL, http.StatusInternalServerError, err)
			ioutils.SendError(w, http.StatusInternalServerError, "internal")
			return
		}
		switch user.Role {
		case models.ParentRole:
			createdDirPath, status, err := fm.userUseCase.CreateParentDirPath(ctx, userRoleID, hashedPathName)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed update parent dir path [status=%d] [error=%s]", r.URL, status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			userDirPath = createdDirPath
		case models.ChildRole:
			createdDirPath, status, err := fm.userUseCase.CreateChildDirPath(ctx, userRoleID, hashedPathName)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed update child dir path [status=%d] [error=%s]", r.URL, status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			userDirPath = createdDirPath
		default:
			fm.logger.Errorf("%s unknown role while creating dir path %s [status=%d]", r.URL, user.Role.String(), http.StatusInternalServerError)
			ioutils.SendError(w, http.StatusInternalServerError, "internal")
			return
		}
	}

	// create a new file in the uploads directory
	fullFilePath := fmt.Sprintf(fm.rootPath + userDirPath + fileHeader.Filename)
	dst, err := os.Create(fullFilePath)
	if err != nil {
		fm.logger.Errorf("%s failed to create new file %s [status=%d] [error=%s]", r.URL, fullFilePath, http.StatusInternalServerError, err)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer dst.Close()

	// copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, bytes.NewReader(buf))
	if err != nil {
		fm.logger.Errorf("%s failed to save user file %s [status=%d] [error=%s]", r.URL, fullFilePath, http.StatusInternalServerError, err)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}
}
