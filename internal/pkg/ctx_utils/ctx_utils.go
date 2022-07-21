package ctx_utils

import (
	"context"

	"github.com/VoyakinH/lokle_backend/internal/models"
)

type CtxKey int

const (
	CtxUser CtxKey = iota
	CtxParent
	CtxChild
	CtxManager
	CtxAdmin
)

func GetUser(ctx context.Context) *models.User {
	user, ok := ctx.Value(CtxUser).(*models.User)
	if ok {
		return user
	}
	return nil
}

func GetParent(ctx context.Context) *models.Parent {
	parent, ok := ctx.Value(CtxParent).(*models.Parent)
	if ok {
		return parent
	}
	return nil
}

func GetChild(ctx context.Context) *models.Child {
	child, ok := ctx.Value(CtxChild).(*models.Child)
	if ok {
		return child
	}
	return nil
}
