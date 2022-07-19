package middleware

import (
	"net/http"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/ioutils"
	"github.com/sirupsen/logrus"
)

func CheckParent(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value(CtxUser).(*models.User)
		if !ok {
			logrus.Errorf("%s failed get ctx user for chekc role with [status=%d]", r.URL, http.StatusForbidden)
			ioutils.SendError(w, http.StatusForbidden, "no auth")
			return
		}
		if user.Role != models.ParentRole {
			logrus.Errorf("%s role %s haven't access to parent functions [status=%d]", r.URL, user.Role, http.StatusForbidden)
			ioutils.SendError(w, http.StatusForbidden, "no auth role")
			return
		}
		h.ServeHTTP(w, r)
	})
}

func CheckChild(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value(CtxUser).(*models.User)
		if !ok {
			logrus.Errorf("%s failed get ctx user for chekc role with [status=%d]", r.URL, http.StatusForbidden)
			ioutils.SendError(w, http.StatusForbidden, "no auth")
			return
		}
		if user.Role != models.ChildRole {
			logrus.Errorf("%s role %s haven't access to child functions [status=%d]", r.URL, user.Role, http.StatusForbidden)
			ioutils.SendError(w, http.StatusForbidden, "no auth role")
			return
		}
		h.ServeHTTP(w, r)
	})
}

func CheckManager(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value(CtxUser).(*models.User)
		if !ok {
			logrus.Errorf("%s failed get ctx user for chekc role with [status=%d]", r.URL, http.StatusForbidden)
			ioutils.SendError(w, http.StatusForbidden, "no auth")
			return
		}
		if user.Role != models.ManagerRole {
			logrus.Errorf("%s role %s haven't access to manager functions [status=%d]", r.URL, user.Role, http.StatusForbidden)
			ioutils.SendError(w, http.StatusForbidden, "no auth role")
			return
		}
		h.ServeHTTP(w, r)
	})
}

func CheckAdmin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, ok := ctx.Value(CtxUser).(*models.User)
		if !ok {
			logrus.Errorf("%s failed get ctx user for chekc role with [status=%d]", r.URL, http.StatusForbidden)
			ioutils.SendError(w, http.StatusForbidden, "no auth")
			return
		}
		if user.Role != models.AdminRole {
			logrus.Errorf("%s role %s haven't access to admin functions [status=%d]", r.URL, user.Role, http.StatusForbidden)
			ioutils.SendError(w, http.StatusForbidden, "no auth role")
			return
		}
		h.ServeHTTP(w, r)
	})
}