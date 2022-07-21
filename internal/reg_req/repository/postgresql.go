package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/VoyakinH/lokle_backend/config"
	"github.com/VoyakinH/lokle_backend/internal/models"
	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

type IPostgresqlRepository interface {
	CreateRegReq(context.Context, uint64, models.RegReqType) (models.RegReqFull, error)
	GetRegRequestList(context.Context, uint64) ([]models.RegReqFull, error)
	GetRegRequestByID(context.Context, uint64) (models.RegReqFull, error)
	DeleteRegReq(context.Context, uint64) (models.RegReqFull, error)
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

func (pr *postgresqlRepository) CreateRegReq(ctx context.Context, uid uint64, reqType models.RegReqType) (models.RegReqFull, error) {
	var req models.RegReqFull
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
		return models.RegReqFull{}, err
	}
	return req, nil
}

func (pr *postgresqlRepository) GetRegRequestList(ctx context.Context, uid uint64) ([]models.RegReqFull, error) {
	rows, err := pr.conn.Query(
		"SELECT id, user_id, type, status, create_time, message FROM registration_requests WHERE user_id = $1;",
		uid,
	)
	if err != nil {
		return []models.RegReqFull{}, err
	}
	defer rows.Close()

	var respList []models.RegReqFull
	var resp models.RegReqFull
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
			return []models.RegReqFull{}, err
		}
		respList = append(respList, resp)
	}
	if err := rows.Err(); err != nil {
		return []models.RegReqFull{}, err
	}
	return respList, nil
}

func (pr *postgresqlRepository) GetRegRequestByID(ctx context.Context, reqID uint64) (models.RegReqFull, error) {
	var req models.RegReqFull
	err := pr.conn.QueryRow(
		`SELECT id, user_id, type, status, create_time, message
		FROM registration_requests
		WHERE id = $1;`,
		reqID,
	).Scan(
		&req.ID,
		&req.UserID,
		&req.Type,
		&req.Status,
		&req.CreateTime,
		&req.Message,
	)
	if err != nil {
		return models.RegReqFull{}, err
	}
	return req, nil
}

func (pr *postgresqlRepository) DeleteRegReq(ctx context.Context, reqID uint64) (models.RegReqFull, error) {
	var deletedReq models.RegReqFull
	err := pr.conn.QueryRow(
		`DELETE FROM registration_requests WHERE id = $1
		RETURNING id, user_id, type, status, create_time, message;`,
		reqID,
	).Scan(
		&deletedReq.ID,
		&deletedReq.UserID,
		&deletedReq.Type,
		&deletedReq.Status,
		&deletedReq.CreateTime,
		&deletedReq.Message,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.RegReqFull{}, nil
		} else {
			return models.RegReqFull{}, err
		}
	}
	return deletedReq, nil
}
