package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/hasher"
	"github.com/VoyakinH/lokle_backend/internal/pkg/mailer"
	pswdgenerator "github.com/VoyakinH/lokle_backend/internal/pkg/psw_generator"
	"github.com/VoyakinH/lokle_backend/internal/pkg/tools"
	"github.com/VoyakinH/lokle_backend/internal/reg_req/repository"
	user_repository "github.com/VoyakinH/lokle_backend/internal/user/repository"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type IRegReqUsecase interface {
	CreateVerifyParentPassportReq(context.Context, models.Parent, models.ParentPassportReq) (int, error)
	GetRegRequestsList(context.Context, uint64) ([]models.RegReqFull, int, error)
	CreateChild(context.Context, models.ChildFirstRegReq, uint64) (models.Child, int, error)
	CompleteRegReq(context.Context, uint64) (int, error)
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
		if existsReq.Type == models.ParentPassportVerification && existsReq.Status == "pending" {
			return http.StatusConflict, fmt.Errorf("RegReqUsecase.CreateVerifyParentPassportReq: parent has already created this request")
		}
	}

	hashedPassport, err := hasher.HashAndSalt(req.Passport)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateVerifyParentPassportReq: failed to hash passport with err: %s", err)
	}
	req.Passport = hashedPassport

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

func (rru *regReqUsecase) GetRegRequestsList(ctx context.Context, uid uint64) ([]models.RegReqFull, int, error) {
	respList, err := rru.psql.GetRegRequestList(ctx, uid)
	if err != nil {
		return []models.RegReqFull{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.GetParentRegRequests: failed to get parent requests with err: %s", err)
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

// func (rru *regReqUsecase) SecondRegistrationChildStage(ctx context.Context, childReq models.ChildSecondRegReq, pid uint64) (models.Child, int, error) {
// 	// rru.userPsql.GetChildByID()
// 	// update child data: passport, places
// 	// rru.userPsql.UpdateChild(ctx, )
// 	// update retationships

// 	// create second type req
// }

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
