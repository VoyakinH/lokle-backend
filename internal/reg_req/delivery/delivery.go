package delivery

import (
	"net/http"
	"strconv"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ctx_utils"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ioutils"
	"github.com/VoyakinH/lokle_backend/internal/pkg/middleware"
	"github.com/VoyakinH/lokle_backend/internal/pkg/tools"
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

	regReqChildAPI := router.PathPrefix("/api/v1/reg/request/child").Subrouter()
	regReqChildAPI.Use(middleware.WithJSON)
	regReqChildAPI.Use(auth.WithAuth)

	regReqChildAPI.Handle("/stage/first", roleMw.CheckParent(http.HandlerFunc(regReqDelivery.FirstSignupChild))).Methods(http.MethodPost)
	regReqChildAPI.Handle("/stage/second", roleMw.CheckParent(http.HandlerFunc(regReqDelivery.SecondSignupChild))).Methods(http.MethodPost)

	regReqCompleteAPI := router.PathPrefix("/api/v1/reg/request").Subrouter()
	regReqCompleteAPI.Use(middleware.WithJSON)
	regReqCompleteAPI.Use(auth.WithAuth)
	regReqCompleteAPI.Use(roleMw.CheckManager)

	regReqCompleteAPI.HandleFunc("/complete", regReqDelivery.CompleteRegReq).Methods(http.MethodGet)
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

	reqList, status, err := rrd.regReqUseCase.GetRegRequestsList(ctx, parent.UserID)
	if err != nil || status != http.StatusOK {
		rrd.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	ioutils.Send(w, status, tools.FullRegReqToSimpleRespList(reqList))
}

func (rrd *RegReqDelivery) FirstSignupChild(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parent := ctx_utils.GetParent(ctx)
	if parent == nil {
		rrd.logger.Errorf("%s failed get ctx parent with [status=%d]", r.URL, http.StatusForbidden)
		ioutils.SendError(w, http.StatusForbidden, "no auth")
		return
	}

	var childReq models.ChildFirstRegReq
	err := ioutils.ReadJSON(r, &childReq)
	if err != nil || !parent.PassportVerified && !childReq.IsStudent {
		rrd.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, err)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	createdChild, status, err := rrd.regReqUseCase.CreateChild(ctx, childReq, parent.ID)
	if err != nil || status != http.StatusOK {
		rrd.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	ioutils.Send(w, status, tools.ChildToChildFullRes(createdChild))
}

func (rrd *RegReqDelivery) SecondSignupChild(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	parent := ctx_utils.GetParent(ctx)
	if parent == nil {
		rrd.logger.Errorf("%s failed get ctx parent with [status=%d]", r.URL, http.StatusForbidden)
		ioutils.SendError(w, http.StatusForbidden, "no auth")
		return
	}

	var childReq models.ChildSecondRegReq
	err := ioutils.ReadJSON(r, &childReq)
	if err != nil {
		rrd.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, err)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	createdReq, status, err := rrd.regReqUseCase.SecondRegistrationChildStage(ctx, childReq, *parent)
	if err != nil || status != http.StatusOK {
		rrd.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	ioutils.Send(w, status, tools.FullRegReqToSimpleResp(createdReq))
}

func (rrd *RegReqDelivery) CompleteRegReq(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqIDString := r.URL.Query().Get("req")
	if reqIDString == "" {
		rrd.logger.Errorf("%s empty query [status=%d]", r.URL, http.StatusBadRequest)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}
	reqID, err := strconv.ParseUint(reqIDString, 10, 64)
	if err != nil {
		rrd.logger.Errorf("%s invalid req id parametr [status=%d]", r.URL, http.StatusBadRequest)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	status, err := rrd.regReqUseCase.CompleteRegReq(ctx, reqID)
	if err != nil || status != http.StatusOK {
		rrd.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	ioutils.SendWithoutBody(w, status)
}
