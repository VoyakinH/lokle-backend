package ctx_utils

import (
	"context"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/middleware"
)

func GetUser(ctx context.Context) *models.User {
	user, ok := ctx.Value(middleware.CtxUser).(*models.User)
	if ok {
		return user
	}
	return nil
}
