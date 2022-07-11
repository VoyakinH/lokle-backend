package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/VoyakinH/lokle_backend/internal/pkg/hasher"
	"github.com/VoyakinH/lokle_backend/internal/pkg/mailer"
	"github.com/VoyakinH/lokle_backend/internal/user/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type IUserUsecase interface {
	CreateSession(context.Context, string, time.Duration) (string, int, error)
	DeleteSession(context.Context, string) (int, error)
	CheckSession(context.Context, string, time.Duration) (string, int, error)
	CreateParent(context.Context, models.Parent) (models.Parent, int, error)
}

type userUsecase struct {
	psql       repository.IPostgresqlRepository
	rdsSession repository.IRedisSessionRepository
	rdsUser    repository.IRedisUserRepository
	logger     logrus.Logger
}

func NewUserUsecase(pr repository.IPostgresqlRepository,
	rsr repository.IRedisSessionRepository,
	rur repository.IRedisUserRepository,
	logger logrus.Logger) IUserUsecase {
	return &userUsecase{
		psql:       pr,
		rdsSession: rsr,
		rdsUser:    rur,
		logger:     logger,
	}
}

const expVerifiedTokenTime = 604800

func (uu *userUsecase) CreateSession(ctx context.Context, email string, sessionExpire time.Duration) (string, int, error) {
	sessionID, err := uuid.NewRandom()
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateSession: failed to create session id with err: %s", err)
	}

	err = uu.rdsSession.CreateSession(ctx, sessionID.String(), email, sessionExpire)
	if err != nil {
		// был 523 код почему?
		return "", http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateSession: failed to create session in redis with err: %s", err)
	}
	return sessionID.String(), http.StatusOK, nil
}

func (uu *userUsecase) DeleteSession(ctx context.Context, cookie string) (int, error) {
	err := uu.rdsSession.DeleteSession(ctx, cookie)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("UserUsecase.DeleteSession: failed to delete session from redis")
	}
	return http.StatusOK, nil
}

func (uu *userUsecase) CheckSession(ctx context.Context, cookie string, expCookieTime time.Duration) (string, int, error) {
	userEmail, err := uu.rdsSession.CheckSession(ctx, cookie)
	if err != nil {
		return "", http.StatusForbidden, fmt.Errorf("UserUsecase.CheckSession: failed to check session in redis")
	}
	if userEmail == "" {
		return "", http.StatusForbidden, fmt.Errorf("UserUsercase.CheckSession: user not authorized")
	}
	err = uu.rdsSession.ProlongSession(ctx, cookie, expCookieTime)
	if err != nil {
		return userEmail, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CheckSession: failed to prolong session in redis")
	}
	return userEmail, http.StatusOK, nil
}

// func (uu *userUsecase) CreateParent(ctx context.Context, parent models.Parent) (models.Parent, int, error) {

// 	return createdParent, http.StatusOK, nil
// }

func (uu *userUsecase) CreateParent(ctx context.Context, parent models.Parent) (models.Parent, int, error) {
	_, err := uu.psql.GetParentByEmail(ctx, parent.Email)
	if err == nil {
		return models.Parent{}, http.StatusConflict, fmt.Errorf("UserUsecase.CreateParent: parent with same email already exists")
	} else if err != nil && err != pgx.ErrNoRows {
		return models.Parent{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParent: failed to check email in db with err: %s", err)
	}

	hashedPswd, err := hasher.HashAndSalt(parent.Password)
	if err != nil {
		return models.Parent{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParent: failed to hash password with err: %s", err)
	}
	parent.Password = hashedPswd

	createdParent, err := uu.psql.CreateParent(ctx, parent)
	if err != nil {
		return models.Parent{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParent: failed to create parent with err: %s", err)
	}

	token, err := uuid.NewRandom()
	if err != nil {
		return models.Parent{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParent: failed to generate token for verification email: %s", err)
	}

	err = uu.rdsUser.AddUserToken(ctx, token.String(), createdParent.Email, expVerifiedTokenTime)
	if err != nil {
		return models.Parent{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParent: failed to save token for verification email to redis: %s", err)
	}

	err = mailer.SendVerifiedEmail(createdParent.Email, token.String())
	if err != nil {
		return models.Parent{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParent: failed to send verification email to parent with err: %s", err)
	}

	createdParent.Password = ""

	return createdParent, http.StatusOK, nil
}
