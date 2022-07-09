package delivery

import (
	"net/http"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ioutils"
	"github.com/VoyakinH/lokle_backend/internal/user/usecase"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type UserDelivery struct {
	UserUseCase usecase.IUserUsecase
	logger      logrus.Logger
}

func SetUserRouting(router *mux.Router, uu usecase.IUserUsecase, logger logrus.Logger) {
	userDelivery := &UserDelivery{
		UserUseCase: uu,
		logger:      logger,
	}

	router.HandleFunc("/api/v1/user/login", userDelivery.CreateUserSession).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/v1/user/logout", userDelivery.DeleteUserSession).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/v1/user/validate_user", userDelivery.CheckUserSession).Methods("GET", "OPTIONS")
}

const expCookieTime = 1382400

func (ud *UserDelivery) CreateUserSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx := r.Context()
	var credentials models.Credentials
	err := ioutils.ReadJSON(r, &credentials)
	if err != nil {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusBadRequest, err)
		ioutils.SendError(w, http.StatusBadRequest, "bad request")
		return
	}

	sessionID, status, err := ud.UserUseCase.CreateSession(ctx, credentials, expCookieTime)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
	}

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  sessionID,
		MaxAge: expCookieTime,
	}

	http.SetCookie(w, cookie)
}

func (ud *UserDelivery) DeleteUserSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx := r.Context()
	cookieToken, err := r.Cookie("session-id")
	if err != nil {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusUnauthorized, err)
		ioutils.SendError(w, http.StatusUnauthorized, "no credentials")
		return
	}

	status, err := ud.UserUseCase.DeleteSession(ctx, cookieToken.Value)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
	}

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

func (ud *UserDelivery) CheckUserSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx := r.Context()
	cookieToken, err := r.Cookie("session-id")
	if err != nil {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, http.StatusUnauthorized, err)
		ioutils.SendError(w, http.StatusUnauthorized, "no credentials")
		return
	}

	status, err := ud.UserUseCase.CheckSession(ctx, cookieToken.Value, expCookieTime)
	if err != nil || status != http.StatusOK {
		ud.logger.Errorf("%s failed with [status=%d] [error=%s]", r.URL, status, err)
		ioutils.SendError(w, status, "internal")
	}

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  cookieToken.Value,
		MaxAge: expCookieTime,
	}

	http.SetCookie(w, cookie)
}
