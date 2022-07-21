package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VoyakinH/lokle_backend/config"
	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type IPostgresqlRepository interface {
	GetUserByEmail(context.Context, string) (models.User, error)
	GetUserByID(context.Context, uint64) (models.User, error)
	GetParentByID(context.Context, uint64) (models.Parent, error)
	GetChildByID(context.Context, uint64) (models.Child, error)
	CreateUser(context.Context, models.User) (models.User, error)
	DeleteUser(context.Context, uint64) (models.User, error)
	VerifyEmail(context.Context, string) (uint64, error)
	CreateParent(context.Context, uint64) (models.Parent, error)
	CreateChild(context.Context, uint64, uint64, models.Child) (models.Child, error)
	CreateParentDirPath(context.Context, uint64, string) (string, error)
	CreateChildDirPath(context.Context, uint64, string) (string, error)
	UpdateParentPassport(context.Context, uint64, string) (string, error)
	VerifyParentPassport(context.Context, uint64) error
	VerifyStageForChild(context.Context, uint64, models.Stage) error
	UpdateUserPswd(context.Context, uint64, string) error
	UpdateChild(context.Context, models.Child) (models.Child, error)
}

type postgresqlRepository struct {
	conn   *pgx.ConnPool
	logger logrus.Logger
}

func NewPostgresqlRepository(cfg config.PostgresConfig, logger logrus.Logger) IPostgresqlRepository {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		cfg.User,
		cfg.DBName,
		cfg.Password,
		cfg.Host,
		cfg.Port)

	pgxConnectionConfig, err := pgx.ParseConnectionString(connStr)
	if err != nil {
		logger.Fatalf("Invalid config string: %s", err)
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxConnectionConfig,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		logger.Fatalf("Error %s occurred during connection to database", err)
	}

	return &postgresqlRepository{conn: pool, logger: logger}
}

func (pr *postgresqlRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := pr.conn.QueryRow(
		"SELECT id, role, first_name, second_name, last_name, phone, email, email_verified, password FROM users WHERE email = $1;",
		email,
	).Scan(
		&user.ID,
		&user.Role,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Phone,
		&user.Email,
		&user.EmailVerified,
		&user.Password,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (pr *postgresqlRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	var createdUser models.User
	err := pr.conn.QueryRow(
		`INSERT INTO users (role, first_name, second_name, last_name, phone, email, password)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, role, first_name, second_name, last_name, phone, email, email_verified;`,
		user.Role,
		user.FirstName,
		user.SecondName,
		user.LastName,
		user.Phone,
		user.Email,
		user.Password,
	).Scan(
		&createdUser.ID,
		&createdUser.Role,
		&createdUser.FirstName,
		&createdUser.SecondName,
		&createdUser.LastName,
		&createdUser.Phone,
		&createdUser.Email,
		&createdUser.EmailVerified,
	)

	if err != nil {
		return models.User{}, err
	}
	return createdUser, nil
}

func (pr *postgresqlRepository) DeleteUser(ctx context.Context, id uint64) (models.User, error) {
	var deletedUser models.User
	err := pr.conn.QueryRow(
		`DELETE FROM users WHERE id = $1
		RETURNING id, first_name, second_name, last_name, phone, email, email_verified;`,
		id,
	).Scan(
		&deletedUser.ID,
		&deletedUser.FirstName,
		&deletedUser.SecondName,
		&deletedUser.LastName,
		&deletedUser.Phone,
		&deletedUser.Email,
		&deletedUser.EmailVerified,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		} else {
			return models.User{}, err
		}
	}
	return deletedUser, nil
}

func (pr *postgresqlRepository) VerifyEmail(ctx context.Context, email string) (uint64, error) {
	var updatedUserID uint64
	err := pr.conn.QueryRow(
		`UPDATE users
		SET email_verified = true
		WHERE email = $1
		RETURNING id;`,
		email,
	).Scan(
		&updatedUserID,
	)

	if err != nil {
		return 0, err
	}
	return updatedUserID, nil
}

func (pr *postgresqlRepository) CreateParent(ctx context.Context, uid uint64) (models.Parent, error) {
	var createdParent models.Parent
	err := pr.conn.QueryRow(
		`INSERT INTO parents (user_id)
		VALUES ($1)
		ON CONFLICT DO NOTHING
		RETURNING id;`,
		uid,
	).Scan(
		&createdParent.ID,
	)

	if err != nil && err != pgx.ErrNoRows {
		return models.Parent{}, err
	}
	return createdParent, nil
}

func (pr *postgresqlRepository) CreateChild(ctx context.Context, uid uint64, pid uint64, child models.Child) (models.Child, error) {
	var createdChild models.Child
	err := pr.conn.QueryRow(
		`INSERT INTO children (user_id, birth_date)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
		RETURNING id, user_id, birth_date;`,
		uid,
		child.BirthDate,
	).Scan(
		&createdChild.ID,
		&createdChild.UserID,
		&createdChild.BirthDate,
	)

	if err != nil && err != pgx.ErrNoRows {
		return models.Child{}, err
	}

	var id uint64
	err = pr.conn.QueryRow(
		`INSERT INTO parents_children (parent_id, child_id)
		VALUES ($1, $2)
		RETURNING id;`,
		pid,
		createdChild.ID,
	).Scan(
		&id,
	)

	if err != nil {
		return models.Child{}, err
	}

	return createdChild, nil
}

func (pr *postgresqlRepository) GetUserByID(ctx context.Context, uid uint64) (models.User, error) {
	var user models.User
	err := pr.conn.QueryRow(
		`SELECT id, role, first_name, second_name, last_name, phone, email, email_verified, password
		FROM users
		WHERE id = $1;`,
		uid,
	).Scan(
		&user.ID,
		&user.Role,
		&user.FirstName,
		&user.SecondName,
		&user.LastName,
		&user.Phone,
		&user.Email,
		&user.EmailVerified,
		&user.Password,
	)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (pr *postgresqlRepository) GetParentByID(ctx context.Context, uid uint64) (models.Parent, error) {
	var parent models.Parent
	err := pr.conn.QueryRow(
		`SELECT
			p.id,
			p.user_id,
			us.role,
			us.first_name,
			us.second_name,
			us.last_name,
			us.phone,
			us.email,
			us.email_verified,
			us.password,
			p.passport,
			p.passport_verified,
			p.dir_path
		FROM users as us
		JOIN parents as p
		ON (p.user_id = us.id)
		WHERE us.id = $1;`,
		uid,
	).Scan(
		&parent.ID,
		&parent.UserID,
		&parent.Role,
		&parent.FirstName,
		&parent.SecondName,
		&parent.LastName,
		&parent.Phone,
		&parent.Email,
		&parent.EmailVerified,
		&parent.Password,
		&parent.Passport,
		&parent.PassportVerified,
		&parent.DirPath,
	)
	if err != nil {
		return models.Parent{}, err
	}
	return parent, nil
}

func (pr *postgresqlRepository) GetChildByID(ctx context.Context, uid uint64) (models.Child, error) {
	var child models.Child
	err := pr.conn.QueryRow(
		`SELECT
			c.id,
			c.user_id,
			us.role,
			us.first_name,
			us.second_name,
			us.last_name,
			us.phone,
			us.email,
			us.email_verified,
			us.password,
			c.birth_date,
			c.done_stage,
			c.passport,
			c.place_of_residence,
			c.place_of_registration,
			c.dir_path
		FROM users as us
		JOIN children as c
		ON (c.user_id = us.id)
		WHERE us.id = $1;`,
		uid,
	).Scan(
		&child.ID,
		&child.UserID,
		&child.Role,
		&child.FirstName,
		&child.SecondName,
		&child.LastName,
		&child.Phone,
		&child.Email,
		&child.EmailVerified,
		&child.Password,
		&child.BirthDate,
		&child.DoneStage,
		&child.Passport,
		&child.PlaceOfResidence,
		&child.PlaceOfRegistration,
		&child.DirPath,
	)
	if err != nil {
		return models.Child{}, err
	}
	return child, nil
}

func (pr *postgresqlRepository) CreateParentDirPath(ctx context.Context, pid uint64, path string) (string, error) {
	var insertedDirPath string
	err := pr.conn.QueryRow(
		`UPDATE parents
		SET dir_path = $2
		WHERE id = $1
		RETURNING dir_path;`,
		pid,
		path,
	).Scan(
		&insertedDirPath,
	)

	if err != nil {
		return "", err
	}
	return insertedDirPath, nil
}

func (pr *postgresqlRepository) CreateChildDirPath(ctx context.Context, cid uint64, path string) (string, error) {
	var insertedDirPath string
	err := pr.conn.QueryRow(
		`UPDATE children
		SET dir_path = $2
		WHERE id = $1
		RETURNING dir_path;`,
		cid,
		path,
	).Scan(
		&insertedDirPath,
	)

	if err != nil {
		return "", err
	}
	return insertedDirPath, nil
}

func (pr *postgresqlRepository) UpdateParentPassport(ctx context.Context, pid uint64, passport string) (string, error) {
	var updatedPassport string
	err := pr.conn.QueryRow(
		`UPDATE parents
		SET passport = $2
		WHERE id = $1
		RETURNING passport;`,
		pid,
		passport,
	).Scan(
		&updatedPassport,
	)

	if err != nil {
		return "", err
	}
	return updatedPassport, nil
}

func (pr *postgresqlRepository) VerifyParentPassport(ctx context.Context, uid uint64) error {
	var updatedPid uint64
	err := pr.conn.QueryRow(
		`UPDATE parents
		SET passport_verified = true
		WHERE user_id = $1
		RETURNING id;`,
		uid,
	).Scan(
		&updatedPid,
	)

	if err != nil {
		return err
	}
	return nil
}

func (pr *postgresqlRepository) VerifyStageForChild(ctx context.Context, uid uint64, completedStage models.Stage) error {
	var updatedCid uint64
	err := pr.conn.QueryRow(
		`UPDATE children
		SET done_stage = $2
		WHERE user_id = $1
		RETURNING id;`,
		uid,
		completedStage,
	).Scan(
		&updatedCid,
	)

	if err != nil {
		return err
	}
	return nil
}

func (pr *postgresqlRepository) UpdateUserPswd(ctx context.Context, uid uint64, newPswd string) error {
	var updatedUid uint64
	err := pr.conn.QueryRow(
		`UPDATE users
		SET password = $2
		WHERE id = $1
		RETURNING id;`,
		uid,
		newPswd,
	).Scan(
		&updatedUid,
	)

	if err != nil {
		return err
	}
	return nil
}

func (pr *postgresqlRepository) UpdateChild(ctx context.Context, child models.Child) (models.Child, error) {
	var updatedChild models.Child
	err := pr.conn.QueryRow(
		`UPDATE children
		SET () = ()
		WHERE id = $1
		RETURNING id;`,
		child.ID,
	).Scan(
		&updatedChild,
	)

	if err != nil {
		return models.Child{}, err
	}
	return updatedChild, nil
}
