package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/VoyakinH/lokle_backend/config"
	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type IPostgresqlRepository interface {
	CreateParentPassportVerification(context.Context, uint64, models.RegReqType) (models.ParentPassportReqFull, error)
	GetParentRegRequestList(context.Context, uint64) ([]models.ParentPassportReqFull, error)
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

func (pr *postgresqlRepository) CreateParentPassportVerification(ctx context.Context, uid uint64, reqType models.RegReqType) (models.ParentPassportReqFull, error) {
	var req models.ParentPassportReqFull
	now := time.Now().Unix()
	err := pr.conn.QueryRow(
		`INSERT INTO registration_requests (user_id, type, create_time)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, type, status, create_time;`,
		uid,
		reqType,
		now,
	).Scan(
		&req.ID,
		&req.UserID,
		&req.Type,
		&req.Status,
		&req.CreateTime,
	)
	if err != nil {
		return models.ParentPassportReqFull{}, err
	}
	return req, nil
}

func (pr *postgresqlRepository) GetParentRegRequestList(ctx context.Context, uid uint64) ([]models.ParentPassportReqFull, error) {
	rows, err := pr.conn.Query(
		"SELECT id, user_id, type, status, create_time, message FROM registration_requests WHERE user_id = $1;",
		uid,
	)
	if err != nil {
		return []models.ParentPassportReqFull{}, err
	}
	defer rows.Close()

	var respList []models.ParentPassportReqFull
	var resp models.ParentPassportReqFull
	for rows.Next() {
		err := rows.Scan(
			&resp.ID,
			&resp.UserID,
			&resp.Type,
			&resp.Status,
			&resp.CreateTime,
			&resp.Message,
		)
		if err != nil {
			return []models.ParentPassportReqFull{}, err
		}
		respList = append(respList, resp)
	}
	if err := rows.Err(); err != nil {
		return []models.ParentPassportReqFull{}, err
	}
	return respList, nil
}
