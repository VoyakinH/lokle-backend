package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/user/repository"
	"github.com/google/uuid"
)

type IUserUsecase interface {
	CreateSession(context.Context, models.Credentials, time.Duration) (string, int, error)
	DeleteSession(context.Context, string) (int, error)
	CheckSession(context.Context, string, time.Duration) (string, int, error)
}

type userUsecase struct {
	Redis repository.IRedisRepository
}

func NewUserUsecase(rr repository.IRedisRepository) IUserUsecase {
	return &userUsecase{
		Redis: rr,
	}
}

func (uu *userUsecase) CreateSession(ctx context.Context, credentials models.Credentials, sessionExpire time.Duration) (string, int, error) {
	sessionID, err := uuid.NewRandom()
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateSession: failed to create session id with err: %s", err)
	}

	err = uu.Redis.CreateSession(ctx, sessionID.String(), credentials.Email, sessionExpire)
	if err != nil {
		// был 523 код почему?
		return "", http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateSession: failed to create session in redis with err: %s", err)
	}
	return sessionID.String(), http.StatusOK, nil
}

func (uu *userUsecase) DeleteSession(ctx context.Context, cookie string) (int, error) {
	err := uu.Redis.DeleteSession(ctx, cookie)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("UserUsecase.DeleteSession: failed to delete session from redis")
	}
	return http.StatusOK, nil
}

func (uu *userUsecase) CheckSession(ctx context.Context, cookie string, expCookieTime time.Duration) (string, int, error) {
	userEmail, err := uu.Redis.CheckSession(ctx, cookie)
	if err != nil {
		return "", http.StatusForbidden, fmt.Errorf("UserUsecase.CheckSession: failed to check session in redis")
	}
	if userEmail == "" {
		return "", http.StatusForbidden, fmt.Errorf("UserUsercase.CheckSession: user not authorized")
	}
	err = uu.Redis.ProlongSession(ctx, cookie, expCookieTime)
	if err != nil {
		return userEmail, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CheckSession: failed to prolong session in redis")
	}
	return userEmail, http.StatusOK, nil
}
