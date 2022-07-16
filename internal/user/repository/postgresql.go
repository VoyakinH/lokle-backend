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
	GetParentByEmail(context.Context, string) (models.Parent, error)
	CreateUserParent(context.Context, models.User) (models.User, error)
	DeleteUser(context.Context, uint64) (models.User, error)
	VerifyEmail(context.Context, string) (uint64, error)
	CreateParent(context.Context, uint64) (models.Parent, error)
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

func (pr *postgresqlRepository) GetParentByEmail(ctx context.Context, email string) (models.Parent, error) {
	var parent models.Parent
	err := pr.conn.QueryRow(
		"SELECT id, first_name, second_name, last_name, phone, email, email_verified FROM parents WHERE email = $1;",
		email,
	).Scan(
		&parent.ID,
		&parent.FirstName,
		&parent.SecondName,
		&parent.LastName,
		&parent.Phone,
		&parent.Email,
		&parent.EmailVerified,
	)
	if err != nil {
		return models.Parent{}, err
	}
	return parent, nil
}

// first registration stage for parents
// now parent haven't passport and other documents
// so we create only default user with role Parent = 0
func (pr *postgresqlRepository) CreateUserParent(ctx context.Context, parent models.User) (models.User, error) {
	var createdParent models.User
	err := pr.conn.QueryRow(
		`INSERT INTO users (role, first_name, second_name, last_name, phone, email, password)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, first_name, second_name, last_name, phone, email, email_verified;`,
		models.ParentRole,
		parent.FirstName,
		parent.SecondName,
		parent.LastName,
		parent.Phone,
		parent.Email,
		parent.Password,
	).Scan(
		&createdParent.ID,
		&createdParent.FirstName,
		&createdParent.SecondName,
		&createdParent.LastName,
		&createdParent.Phone,
		&createdParent.Email,
		&createdParent.EmailVerified,
	)

	if err != nil {
		return models.User{}, err
	}
	return createdParent, nil
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
