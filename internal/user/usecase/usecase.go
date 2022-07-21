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
	CheckSession(context.Context, string) (models.User, int, error)
	ProlongSession(context.Context, string, time.Duration) (int, error)
	CheckUser(context.Context, models.Credentials) (models.User, int, error)
	CreateParentUser(context.Context, models.User) (models.User, int, error)
	VerifyEmail(context.Context, string) (int, error)
	RepeatEmailVerification(context.Context, models.Credentials) (int, error)
	GetUserByID(context.Context, uint64) (models.User, int, error)
	GetParentByID(context.Context, uint64) (models.Parent, int, error)
	GetChildByID(context.Context, uint64) (models.Child, int, error)
	CreateParentDirPath(context.Context, uint64, string) (string, int, error)
	CreateChildDirPath(context.Context, uint64, string) (string, int, error)
	CreateParent(context.Context, uint64) (models.Parent, int, error)
	CreateManager(context.Context, models.User) (models.User, int, error)
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

func (uu *userUsecase) CheckSession(ctx context.Context, cookie string) (models.User, int, error) {
	userEmail, err := uu.rdsSession.CheckSession(ctx, cookie)
	if err != nil {
		return models.User{}, http.StatusForbidden, fmt.Errorf("UserUsecase.CheckSession: failed to check session in redis")
	}
	if userEmail == "" {
		return models.User{}, http.StatusForbidden, fmt.Errorf("UserUsercase.CheckSession: user not authorized")
	}

	user, err := uu.psql.GetUserByEmail(ctx, userEmail)
	if err == pgx.ErrNoRows {
		return models.User{}, http.StatusNotFound, fmt.Errorf("UserUsecase.CheckSession: user with same email not found")
	} else if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CheckSession: failed to check email in db with err: %s", err)
	}

	return user, http.StatusOK, nil
}

func (uu *userUsecase) ProlongSession(ctx context.Context, cookie string, expCookieTime time.Duration) (int, error) {
	err := uu.rdsSession.ProlongSession(ctx, cookie, expCookieTime)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("UserUsecase.ProlongSession: failed to prolong session in redis")
	}
	return http.StatusOK, nil
}

func (uu *userUsecase) checkUserInPSQL(ctx context.Context, credentials models.Credentials) (models.User, int, error) {
	user, err := uu.psql.GetUserByEmail(ctx, credentials.Email)
	if err == pgx.ErrNoRows {
		return models.User{}, http.StatusForbidden, fmt.Errorf("user with same email not found")
	} else if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("failed to check email in db with err: %s", err)
	}

	isAuthorized, err := hasher.ComparePasswords(user.Password, credentials.Password)
	if !isAuthorized || err != nil {
		return models.User{}, http.StatusForbidden, fmt.Errorf("user not authorized: %s", err)
	}

	return user, http.StatusOK, nil
}

func (uu *userUsecase) CheckUser(ctx context.Context, credentials models.Credentials) (models.User, int, error) {
	user, status, err := uu.checkUserInPSQL(ctx, credentials)
	if err != nil || status != http.StatusOK {
		return models.User{}, status, fmt.Errorf("UserUsecase.CheckUser: %s", err)
	}

	return user, http.StatusOK, nil
}

// only for parent
func (uu *userUsecase) createVerificationEmail(ctx context.Context, parent models.User) error {
	token, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate token for verification email: %s", err)
	}

	err = uu.rdsUser.AddUserToken(ctx, token.String(), parent.Email, expVerifiedTokenTime)
	if err != nil {
		return fmt.Errorf("failed to save token for verification email to redis: %s", err)
	}

	err = mailer.SendVerifiedEmail(parent.Email, parent.FirstName, parent.SecondName, token.String())
	if err != nil {
		mailerErr := err
		uu.logger.Infof("delete email verification token for user %s", parent.Email)
		userEmail, err := uu.rdsUser.GetUserAndDelete(ctx, token.String())
		if err != nil {
			uu.logger.Errorf("failed to delete email verification token for user %s", userEmail)
		}
		return fmt.Errorf("failed to send verification email to parent with err: %s", mailerErr)
	}

	return nil
}

func (uu *userUsecase) CreateParentUser(ctx context.Context, parent models.User) (models.User, int, error) {
	_, err := uu.psql.GetUserByEmail(ctx, parent.Email)
	if err == nil {
		return models.User{}, http.StatusConflict, fmt.Errorf("UserUsecase.CreateParentUser: parent with same email already exists")
	} else if err != nil && err != pgx.ErrNoRows {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParentUser: failed to check email in db with err: %s", err)
	}

	hashedPswd, err := hasher.HashAndSalt(parent.Password)
	if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParentUser: failed to hash password with err: %s", err)
	}
	parent.Password = hashedPswd
	parent.Role = models.ParentRole

	createdParent, err := uu.psql.CreateUser(ctx, parent)
	if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParentUser: failed to create parent with err: %s", err)
	}

	if !parent.EmailVerified {
		err = uu.createVerificationEmail(ctx, createdParent)
		if err != nil {
			return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParentUser: %s", err)
		}
	}

	return createdParent, http.StatusOK, nil
}

func (uu *userUsecase) VerifyEmail(ctx context.Context, token string) (int, error) {
	userEmail, err := uu.rdsUser.GetUserAndDelete(ctx, token)
	if err != nil {
		return http.StatusNotFound, fmt.Errorf("UserUsecase.VerifyEmail: failed to get email verification token")
	}
	uid, err := uu.psql.VerifyEmail(ctx, userEmail)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("UserUsecase.VerifyEmail: failed to verify email for user %s with err: %s", userEmail, err)
	}
	_, err = uu.psql.CreateParent(ctx, uid)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("UserUsecase.VerifyEmail: failed to create parent for user %s with err: %s", userEmail, err)
	}
	return http.StatusOK, nil
}

func (uu *userUsecase) RepeatEmailVerification(ctx context.Context, credentials models.Credentials) (int, error) {
	user, status, err := uu.checkUserInPSQL(ctx, credentials)
	if err != nil || status != http.StatusOK {
		return status, fmt.Errorf("UserUsecase.RepeatEmailVerification: %s", err)
	}
	if user.EmailVerified {
		return http.StatusOK, nil
	}

	err = uu.createVerificationEmail(ctx, user)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("UserUsecase.RepeatEmailVerification: %s", err)
	}

	return http.StatusOK, nil
}

func (uu *userUsecase) GetUserByID(ctx context.Context, uid uint64) (models.User, int, error) {
	parent, err := uu.psql.GetUserByID(ctx, uid)
	if err == pgx.ErrNoRows {
		return models.User{}, http.StatusNotFound, fmt.Errorf("UserUsecase.GetUserByID: %s", err)
	} else if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.GetUserByID: user not found")
	}
	return parent, http.StatusOK, nil
}

func (uu *userUsecase) GetParentByID(ctx context.Context, uid uint64) (models.Parent, int, error) {
	parent, err := uu.psql.GetParentByID(ctx, uid)
	if err == pgx.ErrNoRows {
		return models.Parent{}, http.StatusNotFound, fmt.Errorf("UserUsecase.GetParentByID: %s", err)
	} else if err != nil {
		return models.Parent{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.GetParentByID: parent not found")
	}
	return parent, http.StatusOK, nil
}

func (uu *userUsecase) GetChildByID(ctx context.Context, uid uint64) (models.Child, int, error) {
	child, err := uu.psql.GetChildByID(ctx, uid)
	if err == pgx.ErrNoRows {
		return models.Child{}, http.StatusNotFound, fmt.Errorf("UserUsecase.GetChildByID: child not found")
	} else if err != nil {
		return models.Child{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.GetChildByID: %s", err)
	}
	return child, http.StatusOK, nil
}

func (uu *userUsecase) CreateParentDirPath(ctx context.Context, pid uint64, path string) (string, int, error) {
	insertedDirPath, err := uu.psql.CreateParentDirPath(ctx, pid, path)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParentDirPath: %s", err)
	}
	return insertedDirPath, http.StatusOK, nil
}

func (uu *userUsecase) CreateChildDirPath(ctx context.Context, cid uint64, path string) (string, int, error) {
	insertedDirPath, err := uu.psql.CreateChildDirPath(ctx, cid, path)
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateChildDirPath: %s", err)
	}
	return insertedDirPath, http.StatusOK, nil
}

func (uu *userUsecase) CreateParent(ctx context.Context, uid uint64) (models.Parent, int, error) {
	createdParent, err := uu.psql.CreateParent(ctx, uid)
	if err != nil {
		return models.Parent{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateParent: failed to create parent for user %d with err: %s", uid, err)
	}
	return createdParent, http.StatusOK, nil
}

func (uu *userUsecase) CreateManager(ctx context.Context, manager models.User) (models.User, int, error) {
	_, err := uu.psql.GetUserByEmail(ctx, manager.Email)
	if err == nil {
		return models.User{}, http.StatusConflict, fmt.Errorf("UserUsecase.CreateManager: manager with same email already exists")
	} else if err != nil && err != pgx.ErrNoRows {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateManager: failed to check email in db with err: %s", err)
	}

	hashedPswd, err := hasher.HashAndSalt(manager.Password)
	if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateManager: failed to hash password with err: %s", err)
	}
	manager.Password = hashedPswd
	manager.Role = models.ManagerRole

	createdManager, err := uu.psql.CreateUser(ctx, manager)
	if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.CreateManager: failed to create parent with err: %s", err)
	}
	_, err = uu.psql.VerifyEmail(ctx, createdManager.Email)
	if err != nil {
		return models.User{}, http.StatusInternalServerError, fmt.Errorf("UserUsecase.VerCreateManagerifyEmail: failed to verify email for manager %s with err: %s", createdManager.Email, err)
	}
	createdManager.EmailVerified = true

	return createdManager, http.StatusOK, nil
}
