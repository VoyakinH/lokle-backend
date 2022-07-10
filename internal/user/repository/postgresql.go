package repository

import (
	"context"
	"fmt"

	"github.com/VoyakinH/lokle_backend/config"
	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type IPostgresqlRepository interface {
	GetParentByEmail(context.Context, string) (models.Parent, error)
	CreateParent(context.Context, models.Parent) (models.Parent, error)
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

func (pr *postgresqlRepository) CreateParent(ctx context.Context, parent models.Parent) (models.Parent, error) {
	var createdParent models.Parent
	err := pr.conn.QueryRow(
		`INSERT INTO parents (first_name, second_name, last_name, phone, email, password)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, first_name, second_name, last_name, phone, email, email_verified;`,
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
		return models.Parent{}, err
	}
	return createdParent, nil
}