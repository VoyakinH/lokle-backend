package middleware

import (
	"context"
	"net/http"

	"github.com/VoyakinH/lokle_backend/internal/pkg/ioutils"
	"github.com/VoyakinH/lokle_backend/internal/user/usecase"
	"github.com/sirupsen/logrus"
)

type CtxKey int

const (
	CtxUser CtxKey = iota
)

type AuthMiddleware struct {
	UserUseCase usecase.IUserUsecase
	logger      logrus.Logger
}

func NewAuthMiddleware(uu usecase.IUserUsecase, logger logrus.Logger) AuthMiddleware {
	return AuthMiddleware{
		UserUseCase: uu,
		logger:      logger,
	}
}

func (am AuthMiddleware) WithAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cookieToken, err := r.Cookie("session-id")
		if err != nil {
			am.logger.Errorf("%s COOKIE AUTH failed with [status=%d] [error=%s]", r.URL, http.StatusUnauthorized, err)
			ioutils.SendError(w, http.StatusUnauthorized, "no credentials")
			return
		}

		user, status, err := am.UserUseCase.CheckSession(ctx, cookieToken.Value)
		if err != nil || status != http.StatusOK {
			am.logger.Errorf("%s COOKIE AUTH failed with [status=%d] [error=%s]", r.URL, status, err)
			ioutils.SendError(w, status, "internal")
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), CtxUser, &user))

		h.ServeHTTP(w, r)
	})
}

func WithJSON(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		h.ServeHTTP(w, r)
	})
}
