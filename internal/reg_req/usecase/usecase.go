package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/hasher"
	"github.com/VoyakinH/lokle_backend/internal/pkg/tools"
	"github.com/VoyakinH/lokle_backend/internal/reg_req/repository"
	user_repository "github.com/VoyakinH/lokle_backend/internal/user/repository"
	"github.com/sirupsen/logrus"
)

type IRegReqUsecase interface {
	CreateVerifyParentPassportReq(context.Context, models.Parent, models.ParentPassportReq) (int, error)
	GetParentRegRequests(context.Context, uint64) (models.ParentPassportRespList, int, error)
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

	respList, err := rru.psql.GetParentRegRequestList(ctx, parent.UserID)
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

	_, err = rru.psql.CreateParentPassportVerification(ctx, parent.UserID, models.ParentPassportVerification)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.CreateVerifyParentPassportReq: failed to create verification request with err: %s", err)
	}
	return http.StatusOK, nil
}

func (rru *regReqUsecase) GetParentRegRequests(ctx context.Context, uid uint64) (models.ParentPassportRespList, int, error) {
	respList, err := rru.psql.GetParentRegRequestList(ctx, uid)
	if err != nil {
		return models.ParentPassportRespList{}, http.StatusInternalServerError, fmt.Errorf("RegReqUsecase.GetParentRegRequests: failed to get parent requests with err: %s", err)
	}
	return tools.FullParentPassportReqToSimpleList(respList), http.StatusOK, nil
}
