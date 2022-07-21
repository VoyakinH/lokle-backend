package delivery

import (
	"net/http"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ctx_utils"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ioutils"
	"github.com/VoyakinH/lokle_backend/internal/pkg/middleware"
	"github.com/VoyakinH/lokle_backend/internal/reg_req/usecase"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type RegReqDelivery struct {
	regReqUseCase usecase.IRegReqUsecase
	logger        logrus.Logger
}

func SetRegReqRouting(router *mux.Router,
	rru usecase.IRegReqUsecase,
	auth middleware.AuthMiddleware,
	roleMw middleware.RoleMiddleware,
	logger logrus.Logger) {
	regReqDelivery := &RegReqDelivery{
		regReqUseCase: rru,
		logger:        logger,
	}

	regReqParentAPI := router.PathPrefix("/api/v1/reg/request/parent").Subrouter()
	regReqParentAPI.Use(middleware.WithJSON)
	regReqParentAPI.Use(auth.WithAuth)
	regReqParentAPI.Use(roleMw.CheckParent)

	regReqParentAPI.HandleFunc("/passport", regReqDelivery.CreateVerifyParentPassportReq).Methods(http.MethodPost)
	regReqParentAPI.HandleFunc("/list", regReqDelivery.GetParentRegRequests).Methods(http.MethodGet)
}

func (rrd *RegReqDelivery) CreateVerifyParentPassportReq(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parent := ctx_utils.GetParent(ctx)
	if parent == nil {
		rrd.logger.Errorf("%s failed get ctx parent with [status=%d]", r.URL, http.StatusForbidden)
		ioutils.SendError(w, http.StatusForbidden, "no auth")
		return
	}

	var req models.ParentPassportReq
	err := ioutils.ReadJSON(r, &req)
	if err != nil || req.Passport == "" {
		rrd.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, err)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	status, err := rrd.regReqUseCase.CreateVerifyParentPassportReq(ctx, *parent, req)
	if err != nil || status != http.StatusOK {
		rrd.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	ioutils.SendWithoutBody(w, status)
}

func (rrd *RegReqDelivery) GetParentRegRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parent := ctx_utils.GetParent(ctx)
	if parent == nil {
		rrd.logger.Errorf("%s failed get ctx parent with [status=%d]", r.URL, http.StatusForbidden)
		ioutils.SendError(w, http.StatusForbidden, "no auth")
		return
	}

	reqList, status, err := rrd.regReqUseCase.GetParentRegRequests(ctx, parent.UserID)
	if err != nil || status != http.StatusOK {
		rrd.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	ioutils.Send(w, status, reqList)
}
