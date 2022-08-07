package file

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func SetFileRouting(router *mux.Router,
	uu usecase.IUserUsecase,
	auth middleware.AuthMiddleware,
	logger logrus.Logger) FileManager {
	fileManager := FileManager{
		rootPath:    config.File.RootPath,
		userUseCase: uu,
		logger:      logger,
	}

	fileAPI := router.PathPrefix("/api/v1/file/").Subrouter()
	fileAPI.Handle("/upload", auth.WithAuth(http.HandlerFunc(fileManager.Upload))).Methods(http.MethodPost)
	fileAPI.Handle("/download", auth.WithAuth(http.HandlerFunc(fileManager.Download))).Methods(http.MethodPost)
	fileAPI.Handle("/delete", auth.WithAuth(http.HandlerFunc(fileManager.Delete))).Methods(http.MethodPost)

	return fileManager
}

const (
	MaxUploadFilesSize = 120 << 20 // 120MB
	MaxUploadFileSize  = 5 << 20   // 5MB
)

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

	userIDString := r.FormValue("userID")
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		fm.logger.Errorf("%s can't to find user id in req [status=%d]", r.URL, http.StatusBadRequest)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	commonFilename := r.FormValue("filename")
	if err != nil {
		fm.logger.Errorf("%s can't to find filename in req [status=%d]", r.URL, http.StatusBadRequest)
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
		hashedPathName, err := hasher.HashAndSalt(fmt.Sprintf("%s%d", uploadUser.Email, time.Now().Unix()))
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
			createdDirPath, status, err := fm.userUseCase.UpdateParentDirPath(ctx, userID, hashedPathName)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed update parent dir path [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			uploadUser.DirPath = createdDirPath
		case models.ChildRole:
			createdDirPath, status, err := fm.userUseCase.UpdateChildDirPath(ctx, userID, hashedPathName)
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

	sameFilesCount := 0
	err = filepath.Walk(fm.rootPath+uploadUser.DirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fm.logger.Errorf("%s file walk failed with [status=%d] [error=%s]", r.URL, http.StatusInternalServerError, err)
			return nil
		}
		if !info.IsDir() && isEnabledExt(filepath.Ext(path)) && strings.Contains(path, commonFilename) {
			sameFilesCount += 1
		}
		return nil
	})

	if err != nil {
		fm.logger.Errorf("%s failed check files in user dir %s [status=%d]", r.URL, http.StatusInternalServerError)
		ioutils.SendError(w, http.StatusInternalServerError, "internal")
		return
	}

	// validate max file size
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadFilesSize)
	if err := r.ParseMultipartForm(MaxUploadFilesSize); err != nil {
		fm.logger.Errorf("%s files is too big [status=%d]", r.URL, http.StatusBadRequest)
		ioutils.SendError(w, http.StatusBadRequest, "files is too big")
		return
	}

	// find file
	files := r.MultipartForm.File["file"]
	// file, fileHeader, err := r.FormFile("file")
	// if err != nil {
	// 	fm.logger.Errorf("%s can't to find file in req [status=%d]", r.URL, http.StatusBadRequest)
	// 	ioutils.SendError(w, http.StatusBadRequest, "bad request")
	// 	return
	// }
	// defer file.Close()

	for _, fileHeader := range files {
		if fileHeader.Size > MaxUploadFileSize {
			fm.logger.Errorf("%s file is too big [status=%d]", r.URL, http.StatusBadRequest)
			ioutils.SendError(w, http.StatusBadRequest, "file is too big")
			return
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// validate file type
		buf := make([]byte, fileHeader.Size)
		file.Read(buf)
		fileType := http.DetectContentType(buf)
		if !isEnabledFileType(fileType) {
			fm.logger.Errorf("%s forbidden file type %s [status=%d]", r.URL, fileType, http.StatusBadRequest)
			ioutils.SendError(w, http.StatusBadRequest, "bad request")
			return
		}

		// determine filename
		curFilename := fmt.Sprintf("%s_%d%s", commonFilename, sameFilesCount, filepath.Ext(fileHeader.Filename))
		// curFilename := commonFilename + "_" + strconv.Itoa(sameFilesCount) +

		// create a new file in the uploads directory
		fullFilePath := fmt.Sprintf(fm.rootPath + uploadUser.DirPath + curFilename)
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

		// update files with same name count
		sameFilesCount += 1
	}

}

const (
	applicationForAdmissionFileName = "application_for_admission"
	applicationForAdmissionExt      = ".pdf"
	staticFilesFolder               = "$2a$04$pjXLPOhYaTaojItSmcCOc.1z6rzr9pXSWLrBtNLlljzfvTCZGyJA6/"
)

func (fm *FileManager) sendFile(w http.ResponseWriter, filePaths []string, handlerURL string) {
	var resp models.DonwloadResp
	for _, filePath := range filePaths {
		// Read the entire file into a byte slice
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		// Determine the content type of the file
		mimeType := http.DetectContentType(bytes)

		// Append the base64 encoded output
		base64Encoding := base64.StdEncoding.EncodeToString(bytes)
		resp.Files = append(resp.Files, models.FileStruct{
			File: base64Encoding,
			Type: mimeType,
		})
	}
	ioutils.Send(w, http.StatusOK, resp)
}

func (fm *FileManager) Download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx_utils.GetUser(ctx)
	if user == nil {
		fm.logger.Errorf("%s failed get ctx user with [status=%d]", r.URL, http.StatusForbidden)
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
		fm.sendFile(w, []string{filePath}, r.URL.String())
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

	var userFiles []string
	sameFilesCount := 0
	err = filepath.Walk(fm.rootPath+userDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fm.logger.Errorf("%s file walk failed with [status=%d] [error=%s]", r.URL, http.StatusInternalServerError, err)
			return nil
		}
		if !info.IsDir() && isEnabledExt(filepath.Ext(path)) && strings.Contains(path, req.FileName) {
			userFiles = append(userFiles, fmt.Sprintf("%s/%s%s_%d%s", fm.rootPath, userDirPath, req.FileName, sameFilesCount, filepath.Ext(path)))
			sameFilesCount += 1
		}
		return nil
	})

	if err != nil || len(userFiles) == 0 {
		fm.logger.Errorf("%s failed finding file %s [status=%d]", r.URL, req.FileName, http.StatusNotFound)
		ioutils.SendError(w, http.StatusNotFound, "not found")
		return
	}
	// filePath := fmt.Sprintf("%s/%s%s", fm.rootPath, userDirPath, userFile)
	fm.sendFile(w, userFiles, r.URL.String())
}

func (fm *FileManager) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx_utils.GetUser(ctx)
	if user == nil {
		fm.logger.Errorf("%s failed get ctx user with [status=%d]", r.URL, http.StatusForbidden)
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

	var userRole models.Role
	// find user dir path
	switch user.Role {
	case models.ParentRole:
		parent, status, err := fm.userUseCase.GetParentByUID(ctx, user.ID)
		if err != nil || status != http.StatusOK {
			fm.logger.Errorf("%s failed get parent user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
			ioutils.SendError(w, status, "internal")
			return
		}
		if user.ID == req.UserID {
			userRole = parent.Role
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
			userRole = child.Role
		}
	case models.ChildRole:
		if user.ID == req.UserID {
			child, status, err := fm.userUseCase.GetChildByUID(ctx, user.ID)
			if err != nil || status != http.StatusOK {
				fm.logger.Errorf("%s failed get child user [role=%s] [status=%d] [error=%s]", r.URL, user.Role.String(), status, err)
				ioutils.SendError(w, status, "internal")
				return
			}
			userRole = child.Role
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

	err = fm.DeleteFile(ctx, req.UserID, userRole, req.FileName)
	if err != nil {
		fm.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, err)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}
	ioutils.SendWithoutBody(w, http.StatusOK)
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

func (fm *FileManager) DeleteFile(ctx context.Context, uid uint64, userRole models.Role, fileName string) error {
	var userDirPath string
	switch userRole {
	case models.ParentRole:
		parent, status, err := fm.userUseCase.GetParentByUID(ctx, uid)
		if err != nil || status != http.StatusOK {
			return fmt.Errorf("FileManager.DeleteFile: failed get parent user [role=%s] [error=%s]", userRole.String(), err)
		}
		userDirPath = parent.DirPath
	case models.ChildRole:
		child, status, err := fm.userUseCase.GetChildByUID(ctx, uid)
		if err != nil || status != http.StatusOK {
			return fmt.Errorf("FileManager.DeleteFile: failed get child user [role=%s] [error=%s]", userRole.String(), err)
		}
		userDirPath = child.DirPath
	default:
		return fmt.Errorf("FileManager.DeleteFile: unknown role while getting dir path [role=%s]", userRole.String())
	}

	userDir := fmt.Sprintf("%s/%s", fm.rootPath, userDirPath)
	var userFiles []string
	var sameFilesCount int64
	err := filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fm.logger.Errorf("FileManager.DeleteFile: file walk failed with [role=%d] [error=%s]", userRole.String(), err)
			return nil
		}
		if !info.IsDir() && isEnabledExt(filepath.Ext(path)) && strings.Contains(path, fileName) {
			userFiles = append(userFiles, fmt.Sprintf("%s/%s%s_%d%s", fm.rootPath, userDirPath, fileName, sameFilesCount, filepath.Ext(path)))
			sameFilesCount += 1
		}
		return nil
	})

	if err != nil || len(userFiles) == 0 {
		return fmt.Errorf("failed finding file %s [role=%s]", fileName, userRole.String())
	}

	for _, userFile := range userFiles {
		// filePath := fmt.Sprintf("%s/%s", userDir, userFile)
		err = os.Remove(userFile)
		if err != nil {
			return fmt.Errorf("FileManager.DeleteFile: failed to remove user file [role=%s] [error=%s]", userRole.String(), err)
		}

		dirIsEmpty, err := isDirEmpty(userDir)
		if err != nil {
			return fmt.Errorf("FileManager.DeleteFile: failed to check contant of user dir [role=%s] [error=%s]", userRole.String(), err)
		}

		// delete dir path from db
		if dirIsEmpty {
			switch userRole {
			case models.ParentRole:
				_, status, err := fm.userUseCase.UpdateParentDirPath(ctx, uid, "")
				if err != nil || status != http.StatusOK {
					return fmt.Errorf("FileManager.DeleteFile: failed delete parent dir from db [role=%s] [error=%s]", userRole.String(), err)
				}
			case models.ChildRole:
				_, status, err := fm.userUseCase.UpdateChildDirPath(ctx, uid, "")
				if err != nil || status != http.StatusOK {
					return fmt.Errorf("FileManager.DeleteFile: failed delete user dir from db [role=%s] [error=%s]", userRole.String(), err)
				}
			default:
				return fmt.Errorf("FileManager.DeleteFile: unknown role while deleting dir path [role=%s]", userRole.String())
			}

			// rm dir
			err = os.Remove(userDir)
			if err != nil {
				return fmt.Errorf("FileManager.DeleteFile: failed to rm user dir [role=%s] [error=%s]", userRole.String(), err)
			}
		}
	}

	return nil
}

func (fm *FileManager) DeleteDir(ctx context.Context, uid uint64, userRole models.Role) error {
	var userDirPath string
	switch userRole {
	case models.ParentRole:
		parent, status, err := fm.userUseCase.GetParentByUID(ctx, uid)
		if err != nil || status != http.StatusOK {
			return fmt.Errorf("FileManager.DeleteFile: failed get parent user [role=%s] [error=%s]", userRole.String(), err)
		}
		userDirPath = parent.DirPath
	case models.ChildRole:
		child, status, err := fm.userUseCase.GetChildByUID(ctx, uid)
		if err != nil || status != http.StatusOK {
			return fmt.Errorf("FileManager.DeleteFile: failed get child user [role=%s] [error=%s]", userRole.String(), err)
		}
		userDirPath = child.DirPath
	default:
		return fmt.Errorf("FileManager.DeleteFile: unknown role while getting dir path [role=%s]", userRole.String())
	}

	userDir := fmt.Sprintf("%s/%s", fm.rootPath, userDirPath)

	// delete dir path from db
	switch userRole {
	case models.ParentRole:
		_, status, err := fm.userUseCase.UpdateParentDirPath(ctx, uid, "")
		if err != nil || status != http.StatusOK {
			return fmt.Errorf("FileManager.DeleteFile: failed delete parent dir from db [role=%s] [error=%s]", userRole.String(), err)
		}
	case models.ChildRole:
		_, status, err := fm.userUseCase.UpdateChildDirPath(ctx, uid, "")
		if err != nil || status != http.StatusOK {
			return fmt.Errorf("FileManager.DeleteFile: failed delete user dir from db [role=%s] [error=%s]", userRole.String(), err)
		}
	default:
		return fmt.Errorf("FileManager.DeleteFile: unknown role while deleting dir path [role=%s]", userRole.String())
	}

	// rm dir
	err := os.RemoveAll(userDir)
	if err != nil {
		return fmt.Errorf("FileManager.DeleteFile: failed to rm user dir [role=%s] [error=%s]", userRole.String(), err)
	}

	return nil
}
