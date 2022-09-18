package middleware

import (
	"context"
	"net/http"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ctx_utils"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ioutils"
	"github.com/VoyakinH/lokle_backend/internal/user/usecase"
	"github.com/sirupsen/logrus"
)

type RoleMiddleware struct {
	UserUseCase usecase.IUserUsecase
	logger      logrus.Logger
}

func NewRoleMiddleware(uu usecase.IUserUsecase, logger logrus.Logger) RoleMiddleware {
	return RoleMiddleware{
		UserUseCase: uu,
		logger:      logger,
	}
}

func (rm RoleMiddleware) CheckParent(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx_utils.GetUser(ctx)
		if user == nil {
			logrus.Errorf("%s failed get ctx user for check role with [status=%d]", r.URL, http.StatusForbidden)
			ioutils.SendDefaultError(w, http.StatusForbidden)
			return
		}
		if user.Role != models.ParentRole {
			logrus.Errorf("%s role %s haven't access to parent functions [status=%d]", r.URL, user.Role, http.StatusForbidden)
			ioutils.SendDefaultError(w, http.StatusForbidden)
			return
		}

		parent, status, err := rm.UserUseCase.GetParentByUID(ctx, user.ID)
		if err != nil || status != http.StatusOK {
			rm.logger.Errorf("%s get parent from db failed with [status=%d] [error=%s]", r.URL, status, err)
			ioutils.SendDefaultError(w, status)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), ctx_utils.CtxParent, &parent))

		h.ServeHTTP(w, r)
	})
}

func (rm RoleMiddleware) CheckChild(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx_utils.GetUser(ctx)
		if user == nil {
			logrus.Errorf("%s failed get ctx user for check role with [status=%d]", r.URL, http.StatusForbidden)
			ioutils.SendDefaultError(w, http.StatusForbidden)
			return
		}
		if user.Role != models.ChildRole {
			logrus.Errorf("%s role %s haven't access to child functions [status=%d]", r.URL, user.Role, http.StatusForbidden)
			ioutils.SendDefaultError(w, http.StatusForbidden)
			return
		}

		child, status, err := rm.UserUseCase.GetChildByUID(ctx, user.ID)
		if err != nil || status != http.StatusOK {
			rm.logger.Errorf("%s get child from db failed with [status=%d] [error=%s]", r.URL, status, err)
			ioutils.SendDefaultError(w, status)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), ctx_utils.CtxChild, &child))

		h.ServeHTTP(w, r)
	})
}

func (rm RoleMiddleware) CheckManager(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx_utils.GetUser(ctx)
		if user == nil {
			logrus.Errorf("%s failed get ctx user for check role with [status=%d]", r.URL, http.StatusForbidden)
			ioutils.SendDefaultError(w, http.StatusForbidden)
			return
		}
		if user.Role != models.ManagerRole {
			logrus.Errorf("%s role %s haven't access to manager functions [status=%d]", r.URL, user.Role, http.StatusForbidden)
			ioutils.SendDefaultError(w, http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (rm RoleMiddleware) CheckAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := ctx_utils.GetUser(ctx)
		if user == nil {
			logrus.Errorf("%s failed get ctx user for check role with [status=%d]", r.URL, http.StatusForbidden)
			ioutils.SendDefaultError(w, http.StatusForbidden)
			return
		}
		if user.Role != models.AdminRole {
			logrus.Errorf("%s role %s haven't access to admin functions [status=%d]", r.URL, user.Role, http.StatusForbidden)
			ioutils.SendDefaultError(w, http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}
