package file

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	fileAPI.Handle("/download", auth.WithAuth(http.HandlerFunc(fileManager.Download))).Methods(http.MethodPost)
}

const MaxUploadSize = 5 << 20 // 5MB

func isEnabledFileType(fileType string) bool {
	imgTypes := map[string]bool{
		"image/jpg":       true,
		"image/jpeg":      true,
		"image/png":       true,
		"application/pdf": true,
	}
	if imgTypes[fileType] {
		return true
	}
	return false
}

func isEnabledExt(fileType string) bool {
	imgTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".pdf":  true,
	}
	if imgTypes[fileType] {
		return true
	}
	return false
}

type userForUpload struct {
	RoleID  uint64
	Email   string
	Role    models.Role
	DirPath string
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
	if !isEnabledFileType(fileType) {
		fm.logger.Errorf("%s forbidden file type %s [status=%d]", r.URL, fileType, http.StatusBadRequest)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	userIDString := r.FormValue("userID")
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		fm.logger.Errorf("%s can't to find user id in req [status=%d]", r.URL, http.StatusBadRequest)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	// data for user who files are uploaded
	var uploadUser userForUpload
	// find user dir path
	switch user.Role {
	case models.ParentRole:
		parent, status, err := fm.userUseCase.GetParentByUID(ctx, user.ID)
		if err != nil || status != http.StatusOK {
			fm.logger.Errorf("%s failed get parent user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
			ioutils.SendError(w, status, "internal")
			return
		}
		if user.ID == userID {
			uploadUser.DirPath = parent.DirPath
			uploadUser.RoleID = parent.ID
			uploadUser.Email = parent.Email
			uploadUser.Role = parent.Role
		} else {
			child, status, err := fm.userUseCase.GetChildByUID(ctx, userID)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed get child user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			isParentChild, status, err := fm.userUseCase.CheckParentChild(ctx, parent.ID, child.ID)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed check parent's children [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			if !isParentChild {
				fm.logger.Errorf("%s user with userID %d there is not a parent's child [role=%s] [status=%d] [error=%s]", r.URL, userID, user.Role.String(), http.StatusForbidden, err)
				ioutils.SendError(w, http.StatusForbidden, "bad request")
				return
			}
			uploadUser.DirPath = child.DirPath
			uploadUser.RoleID = child.ID
			uploadUser.Email = child.Email
			uploadUser.Role = child.Role
		}
	case models.ChildRole:
		if user.ID == userID {
			child, status, err := fm.userUseCase.GetChildByUID(ctx, user.ID)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed get child user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			uploadUser.DirPath = child.DirPath
			uploadUser.RoleID = child.ID
			uploadUser.Email = child.Email
			uploadUser.Role = child.Role
		} else {
			fm.logger.Errorf("%s child try to load not own files [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), http.StatusForbidden, err)
			ioutils.SendError(w, http.StatusForbidden, "bad request")
			return
		}
	default:
		fm.logger.Errorf("%s unknown role while getting dir path [role=%s] [status=%d]", r.URL, user.Role.String(), http.StatusInternalServerError)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}

	// create user dir path if one not exists
	if uploadUser.DirPath == "" {
		hashedPathName, err := hasher.HashAndSalt(uploadUser.Email)
		if err != nil {
			fm.logger.Errorf("%s failed to create hash for user path with [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), http.StatusInternalServerError, err)
			ioutils.SendError(w, http.StatusInternalServerError, "internal")
			return
		}
		// replace all '/' else we break path
		hashedPathName = strings.ReplaceAll(hashedPathName, "/", "")
		hashedPathName += "/"
		err = os.MkdirAll(fm.rootPath+hashedPathName, os.ModePerm)
		if err != nil {
			fm.logger.Errorf("%s failed to create user's dir with [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), http.StatusInternalServerError, err)
			ioutils.SendError(w, http.StatusInternalServerError, "internal")
			return
		}
		switch uploadUser.Role {
		case models.ParentRole:
			createdDirPath, status, err := fm.userUseCase.CreateParentDirPath(ctx, uploadUser.RoleID, hashedPathName)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed update parent dir path [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			uploadUser.DirPath = createdDirPath
		case models.ChildRole:
			createdDirPath, status, err := fm.userUseCase.CreateChildDirPath(ctx, uploadUser.RoleID, hashedPathName)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed update child dir path [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			uploadUser.DirPath = createdDirPath
		default:
			fm.logger.Errorf("%s unknown role while creating dir path [role=%s] [status=%d]", r.URL, user.Role.String(), http.StatusInternalServerError)
			ioutils.SendError(w, http.StatusInternalServerError, "internal")
			return
		}
	}

	// create a new file in the uploads directory
	fullFilePath := fmt.Sprintf(fm.rootPath + uploadUser.DirPath + fileHeader.Filename)
	dst, err := os.Create(fullFilePath)
	if err != nil {
		fm.logger.Errorf("%s failed to create new file %s [role=%s] [status=%d] [error=%s]", r.URL, fullFilePath, user.Role.String(), http.StatusInternalServerError, err)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer dst.Close()

	// copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, bytes.NewReader(buf))
	if err != nil {
		fm.logger.Errorf("%s failed to save user file %s [role=%s] [status=%d] [error=%s]", r.URL, fullFilePath, user.Role.String(), http.StatusInternalServerError, err)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}
}

const (
	applicationForAdmissionFileName = "application_for_admission"
	applicationForAdmissionExt      = ".jpeg"
	staticFilesFolder               = "$2a$04$pjXLPOhYaTaojItSmcCOc.1z6rzr9pXSWLrBtNLlljzfvTCZGyJA6/"
)

func (fm *FileManager) sendFile(w http.ResponseWriter, filePath string, handlerURL string) {
	file, err := os.Open(filePath)
	if err != nil {
		fm.logger.Errorf("%s failed to open user file %s with [status=%d] [error=%s]", handlerURL, filePath, http.StatusInternalServerError, err)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}
	defer file.Close()

	fileHeader := make([]byte, 512)
	_, err = file.Read(fileHeader)
	if err != nil {
		fm.logger.Errorf("%s failed to read user file %s with [status=%d] [error=%s]", handlerURL, filePath, http.StatusInternalServerError, err)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}
	fileType := http.DetectContentType(fileHeader)

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()

	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")
	w.Header().Set("Content-Disposition", "attachment; filename="+filePath)
	w.Header().Set("Content-Type", fileType)
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))

	file.Seek(0, 0)
	io.Copy(w, file)
}

func (fm *FileManager) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx_utils.GetUser(ctx)
	if user == nil {
		fm.logger.Errorf("%s failed get ctx parent with [status=%d]", r.URL, http.StatusForbidden)
		ioutils.SendError(w, http.StatusForbidden, "no auth")
		return
	}

	var req models.DonwloadReq
	err := ioutils.ReadJSON(r, &req)
	if err != nil || req.FileName == "" {
		fm.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, err)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	// static file
	if req.FileName == applicationForAdmissionFileName {
		filePath := fmt.Sprintf("%s/%s%s%s", fm.rootPath, staticFilesFolder, applicationForAdmissionFileName, applicationForAdmissionExt)
		fm.sendFile(w, filePath, r.URL.String())
		return
	}

	var userDirPath string
	switch user.Role {
	case models.ParentRole:
		parent, status, err := fm.userUseCase.GetParentByUID(ctx, user.ID)
		if err != nil || status != http.StatusOK {
			fm.logger.Errorf("%s failed get parent user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
			ioutils.SendError(w, status, "internal")
			return
		}
		if user.ID == req.UserID {
			userDirPath = parent.DirPath
		} else {
			child, status, err := fm.userUseCase.GetChildByUID(ctx, req.UserID)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed get child user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			isParentChild, status, err := fm.userUseCase.CheckParentChild(ctx, parent.ID, child.ID)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed check parent's children [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			if !isParentChild {
				fm.logger.Errorf("%s user with userID %d there is not a parent's child [role=%s] [status=%d] [error=%s]", r.URL, req.UserID, user.Role.String(), http.StatusForbidden, err)
				ioutils.SendError(w, http.StatusForbidden, "bad request")
				return
			}
			userDirPath = child.DirPath
		}
	case models.ChildRole:
		if user.ID == req.UserID {
			child, status, err := fm.userUseCase.GetChildByUID(ctx, user.ID)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed get child user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			userDirPath = child.DirPath
		} else {
			fm.logger.Errorf("%s child try to download not own files [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), http.StatusForbidden, err)
			ioutils.SendError(w, http.StatusForbidden, "bad request")
			return
		}
	case models.ManagerRole:
		ownerUser, status, err := fm.userUseCase.GetUserByID(ctx, req.UserID)
		if err != nil || status != http.StatusOK {
			fm.logger.Errorf("%s failed get user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
			ioutils.SendError(w, status, "internal")
			return
		}
		if ownerUser.Role == models.ParentRole {
			parent, status, err := fm.userUseCase.GetParentByUID(ctx, ownerUser.ID)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed get parent user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			userDirPath = parent.DirPath
		} else if ownerUser.Role == models.ChildRole {
			child, status, err := fm.userUseCase.GetChildByUID(ctx, ownerUser.ID)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed get child user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			userDirPath = child.DirPath
		} else {
			fm.logger.Errorf("%s manager try to download not parent or child file [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), http.StatusForbidden, err)
			ioutils.SendError(w, http.StatusForbidden, "bad request")
			return
		}
	default:
		fm.logger.Errorf("%s unknown role while getting dir path [role=%s] [status=%d]", r.URL, user.Role.String(), http.StatusInternalServerError)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}

	var userFile string
	err = filepath.Walk(fm.rootPath+userDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fm.logger.Errorf("%s file walk failed with [status=%d] [error=%s]", r.URL, http.StatusInternalServerError, err)
			return nil
		}
		if !info.IsDir() && isEnabledExt(filepath.Ext(path)) && strings.Contains(path, req.FileName) {
			if userFile != "" {
				return fmt.Errorf("found some files with this name")
			}
			userFile = req.FileName + filepath.Ext(path)
		}
		return nil
	})

	if err != nil {
		fm.logger.Errorf("%s failed finding file %s [status=%d]", r.URL, req.FileName, http.StatusInternalServerError)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}
	filePath := fmt.Sprintf("%s/%s%s", fm.rootPath, userDirPath, userFile)
	fm.sendFile(w, filePath, r.URL.String())
}
