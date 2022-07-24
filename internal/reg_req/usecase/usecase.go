package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/crypt"
	"github.com/VoyakinH/lokle_backend/internal/pkg/hasher"
	"github.com/VoyakinH/lokle_backend/internal/pkg/mailer"
	pswdgenerator "github.com/VoyakinH/lokle_backend/internal/pkg/psw_generator"
	"github.com/VoyakinH/lokle_backend/internal/pkg/tools"
	"github.com/VoyakinH/lokle_backend/internal/reg_req/repository"
	user_repository "github.com/VoyakinH/lokle_backend/internal/user/repository"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

const (
	PendingReqStatus = "pending"
	FailedReqStatus  = "failed"
)

type IRegReqUsecase interface {
	CreateVerifyParentPassportReq(context.Context, models.Parent, models.ParentPassportReq) (int, error)
	GetRegRequestsList(context.Context, uint64) ([]models.RegReqFull, int, error)
	GetRegRequestsListAll(context.Context) ([]models.RegReqWithUser, int, error)
	CreateChild(context.Context, models.ChildFirstRegReq, uint64) (models.Child, int, error)
	CompleteRegReq(context.Context, uint64) (int, error)
	SecondRegistrationChildStage(context.Context, models.ChildSecondRegReq, models.Parent) (models.RegReqFull, int, error)
	ThirdRegistrationChildStage(context.Context, models.ChildThirdRegReq) (models.RegReqFull, int, error)
	FailedRegReq(context.Context, uint64, models.FailedReq) (int, error)
	FixVerifyParentPassportReq(context.Context, models.Parent, models.FixParentPassportReq) (int, error)
	FixChild(context.Context, models.FixChildFirstRegReq) (int, error)
	FixSecondRegistrationChildStage(context.Context, models.FixChildSecondRegReq, models.Parent) (int, error)
	FixThirdRegistrationChildStage(context.Context, models.FixChildThirdRegReq) (int, error)
}

type regReqUsecase struct {
	psql     repository.IPostgresqlRepository
	userPsql user_repository.IPostgresqlRepository
	logger   logrus.Logger
}

func NewRegReqUsecase(pr repository.IPostgresqlRepository, ur user_repository.IPostgresqlRepository, logger logrus.Logger) IRegReqUsecase {
	return &regReqUsecase{
		psql:     pr,
		userPsql: ur,
		logger:   logger,
	}
}

func (rru *regReqUsecase) CreateVerifyParentPassportReq(ctx context.Context, parent models.Parent, req models.ParentPassportReq) (int, error) {
	if parent.PassportVerified {
		return http.StatusOK, fmt.Errorf("RegReqUsecase.CreateVerifyParentPassportReq: parent passport has been already verified")
	}

	respList, err := rru.psql.GetRegRequestList(ctx, parent.UserID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateVerifyParentPassportReq: failed to get parent's requests: %s", err)
	}
	for _, existsReq := range respList {
		if existsReq.Type == models.ParentPassportVerification && existsReq.Status == PendingReqStatus {
			return http.StatusConflict, fmt.Errorf("RegReqUsecase.CreateVerifyParentPassportReq: parent has already created this request")
		}
	}

	encryptedPassport, err := crypt.Encrypt(req.Passport)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateVerifyParentPassportReq: failed to encrypt parent passport with err: %s", err)
	}
	req.Passport = encryptedPassport

	_, err = rru.userPsql.UpdateParentPassport(ctx, parent.ID, req.Passport)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateVerifyParentPassportReq: failed to update parent passport with err: %s", err)
	}

	_, err = rru.psql.CreateRegReq(ctx, parent.UserID, models.ParentPassportVerification)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateVerifyParentPassportReq: failed to create verification request with err: %s", err)
	}

	return http.StatusOK, nil
}

func (rru *regReqUsecase) FixVerifyParentPassportReq(ctx context.Context, parent models.Parent, reqFix models.FixParentPassportReq) (int, error) {
	if parent.PassportVerified {
		return http.StatusOK, fmt.Errorf("RegReqUsecase.FixVerifyParentPassportReq: parent passport has been already verified")
	}

	_, err := rru.psql.GetRegRequestByID(ctx, reqFix.ReqID)
	if err == pgx.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.FixVerifyParentPassportReq: request not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixVerifyParentPassportReq: failed to get request with err: %s", err)
	}

	encryptedPassport, err := crypt.Encrypt(reqFix.Passport)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixVerifyParentPassportReq: failed to encrypt parent passport with err: %s", err)
	}
	reqFix.Passport = encryptedPassport

	_, err = rru.userPsql.UpdateParentPassport(ctx, parent.ID, reqFix.Passport)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixVerifyParentPassportReq: failed to update parent passport with err: %s", err)
	}

	err = rru.psql.FixRegReq(ctx, reqFix.ReqID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixVerifyParentPassportReq: failed to fix verification request with err: %s", err)
	}

	return http.StatusOK, nil
}

func (rru *regReqUsecase) GetRegRequestsList(ctx context.Context, uid uint64) ([]models.RegReqFull, int, error) {
	respList, err := rru.psql.GetRegRequestList(ctx, uid)
	if err != nil {
		return []models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.GetParentRegRequests: failed to get parent requests with err: %s", err)
	}
	return respList, http.StatusOK, nil
}

func (rru *regReqUsecase) GetRegRequestsListAll(ctx context.Context) ([]models.RegReqWithUser, int, error) {
	respList, err := rru.psql.GetRegRequestListAll(ctx)
	if err != nil {
		return []models.RegReqWithUser{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.GetRegRequestsListAll: failed to get parent requests with err: %s", err)
	}
	now := uint64(time.Now().Unix())
	for i := range respList {
		respList[i].TimeInQueue = uint32((now - respList[i].CreateTime) / 86400)
	}
	return respList, http.StatusOK, nil
}

func (rru *regReqUsecase) CreateChild(ctx context.Context, childReq models.ChildFirstRegReq, pid uint64) (models.Child, int, error) {
	child := childReq.Child
	_, err := rru.userPsql.GetUserByEmail(ctx, child.Email)
	if err == nil {
		return models.Child{}, http.StatusConflict, fmt.Errorf("RegReqUsecase.CreateChild: child with same email already exists")
	} else if err != nil && err != pgx.ErrNoRows {
		return models.Child{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateChild: failed to check email in db with err: %s", err)
	}
	child.Password = ""
	child.Role = models.ChildRole

	createdChildUser, err := rru.userPsql.CreateUser(ctx, tools.ChildToUser(child))
	if err != nil {
		return models.Child{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateChild: failed to create child user with err: %s", err)
	}

	createdChild, err := rru.userPsql.CreateChild(ctx, createdChildUser.ID, pid, child)
	if err != nil {
		return models.Child{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateChild: failed to create child with err: %s", err)
	}

	var reqType models.RegReqType
	if childReq.IsStudent {
		reqType = models.ChildFirstStageForStudent
	} else {
		reqType = models.ChildFirstStage
	}
	_, err = rru.psql.CreateRegReq(ctx, createdChild.UserID, reqType)
	if err != nil {
		return models.Child{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateChild: failed to create first stage request with err: %s", err)
	}

	return models.Child{
		ID:            createdChild.ID,
		UserID:        createdChild.UserID,
		Role:          createdChildUser.Role,
		FirstName:     createdChildUser.FirstName,
		SecondName:    createdChildUser.SecondName,
		LastName:      createdChildUser.LastName,
		Email:         createdChildUser.Email,
		EmailVerified: createdChild.EmailVerified,
		Phone:         createdChildUser.Phone,
		BirthDate:     createdChild.BirthDate,
	}, http.StatusOK, nil
}

func (rru *regReqUsecase) FixChild(ctx context.Context, childReq models.FixChildFirstRegReq) (int, error) {
	req, err := rru.psql.GetRegRequestByID(ctx, childReq.ReqID)
	if err == pgx.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: request not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: failed to get request with err: %s", err)
	}
	if req.Status == PendingReqStatus {
		return http.StatusConflict, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: can't update request in pending status")
	}

	err = rru.userPsql.UpdateChild(ctx, childReq.Child)
	if err == pgx.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.FixChild: child not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixChild: failed to update child with err: %s", err)
	}

	childUser, err := rru.userPsql.GetUserByID(ctx, childReq.Child.UserID)
	if err == pgx.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.FixChild: user not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixChild: failed to get user with err: %s", err)
	}
	if childUser.Email != childReq.Child.Email {
		err = rru.userPsql.UpdateUserWithEmail(ctx, models.User{
			ID:         childReq.Child.UserID,
			FirstName:  childReq.Child.FirstName,
			SecondName: childReq.Child.SecondName,
			LastName:   childReq.Child.LastName,
			Phone:      childReq.Child.Phone,
			Email:      childReq.Child.Email,
		})
	} else {
		err = rru.userPsql.UpdateUserWithoutEmail(ctx, models.User{
			ID:         childReq.Child.UserID,
			FirstName:  childReq.Child.FirstName,
			SecondName: childReq.Child.SecondName,
			LastName:   childReq.Child.LastName,
			Phone:      childReq.Child.Phone,
		})
	}

	if err == pgx.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.FixChild: user not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixChild: failed to update user with err: %s", err)
	}

	err = rru.psql.FixRegReq(ctx, childReq.ReqID)
	if err == pgx.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.FixChild: request not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixChild: failed to fix request with err: %s", err)
	}

	return http.StatusOK, nil
}

func (rru *regReqUsecase) SecondRegistrationChildStage(ctx context.Context, childReq models.ChildSecondRegReq, parent models.Parent) (models.RegReqFull, int, error) {
	child, err := rru.userPsql.GetChildByUID(ctx, childReq.Child.UserID)
	if err != nil {
		return models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.SecondRegistrationChildStage: failed to get child data with err: %s", err)
	}

	respList, err := rru.psql.GetRegRequestList(ctx, child.UserID)
	if err != nil {
		return models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.SecondRegistrationChildStage: failed to get child's requests: %s", err)
	}
	for _, existsReq := range respList {
		if (existsReq.Type == models.ChildFirstStage ||
			existsReq.Type == models.ChildSecondStage ||
			existsReq.Type == models.ChildThirdStage) &&
			existsReq.Status == PendingReqStatus {
			return models.RegReqFull{}, http.StatusConflict, fmt.Errorf("RegReqUsecase.SecondRegistrationChildStage: child has already created this request")
		}
	}
	childReq.Child.BirthDate = child.BirthDate
	encryptedPassport, err := crypt.Encrypt(childReq.Child.Passport)
	if err != nil {
		return models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.SecondRegistrationChildStage: failed to encrypt child passport with err: %s", err)
	}
	childReq.Child.Passport = encryptedPassport

	err = rru.userPsql.UpdateChild(ctx, childReq.Child)
	if err != nil {
		return models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.SecondRegistrationChildStage: failed to update child data with err: %s", err)
	}

	err = rru.userPsql.UpdateParentChildRelationship(ctx, parent.ID, child.ID, childReq.Relationship)
	if err != nil {
		return models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.SecondRegistrationChildStage: failed to update parent and child relationship with err: %s", err)
	}

	req, err := rru.psql.CreateRegReq(ctx, child.UserID, models.ChildSecondStage)
	if err != nil {
		return models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.SecondRegistrationChildStage: failed to create verification request with err: %s", err)
	}

	return req, http.StatusOK, nil
}

func (rru *regReqUsecase) FixSecondRegistrationChildStage(ctx context.Context, childReq models.FixChildSecondRegReq, parent models.Parent) (int, error) {
	req, err := rru.psql.GetRegRequestByID(ctx, childReq.ReqID)
	if err == pgx.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: request not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: failed to get request with err: %s", err)
	}
	if req.Status == PendingReqStatus {
		return http.StatusConflict, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: can't update request in pending status")
	}

	child, err := rru.userPsql.GetChildByUID(ctx, childReq.Child.UserID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixSecondRegistrationChildStage: failed to get child data with err: %s", err)
	}
	childReq.Child.BirthDate = child.BirthDate
	encryptedPassport, err := crypt.Encrypt(childReq.Child.Passport)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixSecondRegistrationChildStage: failed to encrypt child passport with err: %s", err)
	}
	childReq.Child.Passport = encryptedPassport

	err = rru.userPsql.UpdateChild(ctx, childReq.Child)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixSecondRegistrationChildStage: failed to update child data with err: %s", err)
	}

	err = rru.userPsql.UpdateParentChildRelationship(ctx, parent.ID, child.ID, childReq.Relationship)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixSecondRegistrationChildStage: failed to update parent and child relationship with err: %s", err)
	}

	err = rru.psql.FixRegReq(ctx, childReq.ReqID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixSecondRegistrationChildStage: failed to fix request with err: %s", err)
	}

	return http.StatusOK, nil
}

func (rru *regReqUsecase) ThirdRegistrationChildStage(ctx context.Context, childReq models.ChildThirdRegReq) (models.RegReqFull, int, error) {
	child, err := rru.userPsql.GetChildByUID(ctx, childReq.Child.UserID)
	if err != nil {
		return models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.ThirdRegistrationChildStage: failed to get child data with err: %s", err)
	}

	respList, err := rru.psql.GetRegRequestList(ctx, child.UserID)
	if err != nil {
		return models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.ThirdRegistrationChildStage: failed to get child's requests: %s", err)
	}
	for _, existsReq := range respList {
		if (existsReq.Type == models.ChildFirstStage ||
			existsReq.Type == models.ChildSecondStage ||
			existsReq.Type == models.ChildThirdStage) &&
			existsReq.Status == PendingReqStatus {
			return models.RegReqFull{}, http.StatusConflict, fmt.Errorf("RegReqUsecase.ThirdRegistrationChildStage: child has already created this request")
		}
	}

	req, err := rru.psql.CreateRegReq(ctx, child.UserID, models.ChildThirdStage)
	if err != nil {
		return models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.ThirdRegistrationChildStage: failed to create verification request with err: %s", err)
	}

	return req, http.StatusOK, nil
}

func (rru *regReqUsecase) FixThirdRegistrationChildStage(ctx context.Context, childReq models.FixChildThirdRegReq) (int, error) {
	req, err := rru.psql.GetRegRequestByID(ctx, childReq.ReqID)
	if err == pgx.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: request not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: failed to get request with err: %s", err)
	}
	if req.Status == PendingReqStatus {
		return http.StatusConflict, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: can't update request in pending status")
	}

	err = rru.psql.FixRegReq(ctx, childReq.ReqID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FixThirdRegistrationChildStage: failed to create verification request with err: %s", err)
	}

	return http.StatusOK, nil
}

func (rru *regReqUsecase) completeThirdRegistrationChildStage(ctx context.Context, uid uint64) error {
	err := rru.userPsql.VerifyStageForChild(ctx, uid, models.ThirdStage)
	if err != nil {
		return fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to verify first stage for old students in db with err: %s", err)
	}
	user, err := rru.userPsql.GetUserByID(ctx, uid)
	if err != nil {
		return fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to find user for child with err: %s", err)
	}
	childPswd := pswdgenerator.GeneratePassword(10, 0, 2, 2)
	hashedPswd, err := hasher.HashAndSalt(childPswd)
	if err != nil {
		return fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to hash password with err: %s", err)
	}
	err = rru.userPsql.UpdateUserPswd(ctx, user.ID, hashedPswd)
	if err != nil {
		return fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to update child password with err: %s", err)
	}
	err = mailer.SendCompleteChildRegistrationEmail(user.Email, user.FirstName, user.SecondName, childPswd)
	if err != nil {
		return fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to send email with credentials to child with err: %s", err)
	}
	return nil
}

func (rru *regReqUsecase) CompleteRegReq(ctx context.Context, reqID uint64) (int, error) {
	req, err := rru.psql.GetRegRequestByID(ctx, reqID)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to find request with err: %s", err)
	}
	if req.Status == FailedReqStatus {
		return http.StatusConflict, fmt.Errorf("RegReqUsecase.CompleteRegReq: request has been already in failed status")
	}

	switch req.Type {
	case models.ParentPassportVerification:
		err = rru.userPsql.VerifyParentPassport(ctx, req.UserID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to verify parent passport in db with err: %s", err)
		}
	case models.ChildFirstStageForStudent:
		err = rru.completeThirdRegistrationChildStage(ctx, req.UserID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CompleteRegReq: %s", err)
		}
	case models.ChildFirstStage:
		err = rru.userPsql.VerifyStageForChild(ctx, req.UserID, models.FirstStage)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to verify first stage for new students in db with err: %s", err)
		}
	case models.ChildSecondStage:
		err = rru.userPsql.VerifyStageForChild(ctx, req.UserID, models.SecondStage)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to verify first stage for new students in db with err: %s", err)
		}
	case models.ChildThirdStage:
		err = rru.completeThirdRegistrationChildStage(ctx, req.UserID)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CompleteRegReq: %s", err)
		}
	default:
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CompleteRegReq: unknown request type %s", req.Type.String())
	}

	_, err = rru.psql.DeleteRegReq(ctx, reqID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CompleteRegReq: failed to delete request after completed with err: %s", err)
	}

	return http.StatusOK, nil
}

func (rru *regReqUsecase) FailedRegReq(ctx context.Context, managerID uint64, failedReq models.FailedReq) (int, error) {
	req, err := rru.psql.GetRegRequestByID(ctx, failedReq.ReqId)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FailedRegReq: failed to get request with err: %s", err)
	}
	if req.Status == FailedReqStatus {
		return http.StatusConflict, fmt.Errorf("RegReqUsecase.FailedRegReq: request has been already in failed status")
	}
	err = rru.psql.FailedRegReq(ctx, managerID, failedReq)
	if err == pgx.ErrNoRows {
		return http.StatusNotFound, fmt.Errorf("RegReqUsecase.FailedRegReq: request not found")
	} else if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.FailedRegReq: failed to update request with err: %s", err)
	}

	return http.StatusOK, nil
}
