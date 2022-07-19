package delivery

import (
	"net/http"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ctx_utils"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ioutils"
	"github.com/VoyakinH/lokle_backend/internal/pkg/middleware"
	"github.com/VoyakinH/lokle_backend/internal/pkg/tools"
	"github.com/VoyakinH/lokle_backend/internal/user/usecase"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UserDelivery struct {
	userUseCase usecase.IUserUsecase
	logger      logrus.Logger
}

func SetUserRouting(router *mux.Router, uu usecase.IUserUsecase, auth middleware.AuthMiddleware, logger logrus.Logger) {
	userDelivery := &UserDelivery{
		userUseCase: uu,
		logger:      logger,
	}

	userAPI := router.PathPrefix("/api/v1/user/").Subrouter()
	userAPI.Use(middleware.WithJSON)

	userAPI.HandleFunc("/auth", userDelivery.CreateUserSession).Methods(http.MethodPost)
	userAPI.HandleFunc("/auth", userDelivery.DeleteUserSession).Methods(http.MethodDelete)
	userAPI.Handle("/auth", auth.WithAuth(http.HandlerFunc(userDelivery.CheckUserSession))).Methods(http.MethodGet)

	userAPI.HandleFunc("/parent", userDelivery.SignupParent).Methods(http.MethodPost)
	userAPI.HandleFunc("/parent", userDelivery.GetParent).Methods(http.MethodGet)

	userAPI.HandleFunc("/email", userDelivery.EmailVerification).Methods(http.MethodGet)
	userAPI.HandleFunc("/email", userDelivery.RepeatEmailVerification).Methods(http.MethodPost)
}

const expCookieTime = 1382400

func (ud *UserDelivery) CreateUserSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var credentials models.Credentials
	err := ioutils.ReadJSON(r, &credentials)
	if err != nil || credentials.Email == "" || credentials.Password == "" {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, err)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	user, status, err := ud.userUseCase.CheckUser(ctx, credentials)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "failed login")
		return
	}

	// if email not verified we send user but doesn't create and set cookie
	if !user.EmailVerified {
		ioutils.Send(w, status, user)
		return
	}

	sessionID, status, err := ud.userUseCase.CreateSession(ctx, credentials.Email, expCookieTime)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  sessionID,
		MaxAge: expCookieTime,
		Path:   "/api/v1",
	}

	http.SetCookie(w, cookie)
	ioutils.Send(w, status, tools.UserToUserRes(user))
}

func (ud *UserDelivery) DeleteUserSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cookieToken, err := r.Cookie("session-id")
	if err != nil {
		ud.logger.Warnf("%s cookie not found with [status=%d] [error=%s]", r.URL, http.StatusOK, err)
		ioutils.SendError(w, http.StatusOK, "no credentials")
		return
	}

	status, err := ud.userUseCase.DeleteSession(ctx, cookieToken.Value)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
		Path:   "/api/v1",
	}

	http.SetCookie(w, cookie)
}

func (ud *UserDelivery) CheckUserSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := ctx_utils.GetUser(ctx)
	if user == nil {
		ud.logger.Errorf("%s failed get ctx user with [status=%d]", r.URL, http.StatusForbidden)
		ioutils.SendError(w, http.StatusForbidden, "no auth")
		return
	}

	cookieToken, err := r.Cookie("session-id")
	if err != nil {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusUnauthorized, err)
		ioutils.SendError(w, http.StatusUnauthorized, "no credentials")
		return
	}

	status, err := ud.userUseCase.ProlongSession(ctx, cookieToken.Value, expCookieTime)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  cookieToken.Value,
		MaxAge: expCookieTime,
		Path:   "/api/v1",
	}

	http.SetCookie(w, cookie)
	ioutils.Send(w, status, tools.UserToUserRes(*user))
}

func (ud *UserDelivery) SignupParent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var parent models.User
	err := ioutils.ReadJSON(r, &parent)
	if err != nil || parent.Email == "" || parent.Password == "" {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, err)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	createdParent, status, err := ud.userUseCase.CreateParent(ctx, parent)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	ioutils.Send(w, status, tools.UserToUserRes(createdParent))
}

func (ud *UserDelivery) EmailVerification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.URL.Query().Get("token")
	if token == "" {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, "empty token")
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	status, err := ud.userUseCase.VerifyEmail(ctx, token)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
		return
	}

	ioutils.SendWithoutBody(w, status)
}

func (ud *UserDelivery) RepeatEmailVerification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var credentials models.Credentials
	err := ioutils.ReadJSON(r, &credentials)
	if err != nil || credentials.Email == "" || credentials.Password == "" {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, err)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	status, err := ud.userUseCase.RepeatEmailVerification(ctx, credentials)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "failed login")
		return
	}

	ioutils.SendWithoutBody(w, status)
}

func (ud *UserDelivery) GetParent(w http.ResponseWriter, r *http.Request) {

}
