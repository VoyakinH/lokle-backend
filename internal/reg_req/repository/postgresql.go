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
	FixRegReq(context.Context, uint64) error
	GetRegRequestList(context.Context, uint64) ([]models.RegReqFull, error)
	GetRegRequestListAll(context.Context) ([]models.RegReqWithUser, error)
	GetRegRequestByID(context.Context, uint64) (models.RegReqFull, error)
	DeleteRegReq(context.Context, uint64) (models.RegReqFull, error)
	FailedRegReq(context.Context, uint64, models.FailedReq) error
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

type managerNull struct {
	ID         *uint64
	FirstName  *string
	SecondName *string
	LastName   *string
	Role       *models.Role
}

func (mn *managerNull) convertToManager() *models.User {
	result := &models.User{}
	isEmpty := true
	if mn.ID != nil {
		result.ID = *mn.ID
		isEmpty = false
	}
	if mn.FirstName != nil {
		result.FirstName = *mn.FirstName
		isEmpty = false
	}
	if mn.SecondName != nil {
		result.SecondName = *mn.SecondName
		isEmpty = false
	}
	if mn.LastName != nil {
		result.LastName = *mn.LastName
		isEmpty = false
	}
	if mn.Role != nil {
		result.Role = *mn.Role
		isEmpty = false
	}
	if isEmpty {
		return nil
	}
	return result
}

func (pr *postgresqlRepository) GetRegRequestListAll(ctx context.Context) ([]models.RegReqWithUser, error) {
	rows, err := pr.conn.Query(
		`SELECT
			rr.id,
			us.id,
			us.first_name,
			us.second_name,
			us.last_name,
			us.role,
			us.email,
			us.phone,
			usm.id,
			usm.first_name,
			usm.second_name,
			usm.last_name,
			usm.role,
			rr.type,
			rr.status,
			rr.create_time,
			rr.message
		FROM registration_requests AS rr
		JOIN users AS us ON (us.id = rr.user_id)
		LEFT JOIN users AS usm ON (usm.id = rr.manager_id)
		WHERE rr.status = 'pending'
		ORDER BY create_time;`,
	)
	if err != nil {
		return []models.RegReqWithUser{}, err
	}
	defer rows.Close()

	var respList []models.RegReqWithUser
	var resp models.RegReqWithUser
	var tempManager managerNull
	for rows.Next() {
		err := rows.Scan(
			&resp.ID,
			&resp.User.ID,
			&resp.User.FirstName,
			&resp.User.SecondName,
			&resp.User.LastName,
			&resp.User.Role,
			&resp.User.Email,
			&resp.User.Phone,
			&tempManager.ID,
			&tempManager.FirstName,
			&tempManager.SecondName,
			&tempManager.LastName,
			&tempManager.Role,
			&resp.Type,
			&resp.Status,
			&resp.CreateTime,
			&resp.Message,
		)
		if err != nil {
			return []models.RegReqWithUser{}, err
		}
		resp.Manager = tempManager.convertToManager()
		respList = append(respList, resp)
	}
	if err := rows.Err(); err != nil {
		return []models.RegReqWithUser{}, err
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

func (pr *postgresqlRepository) FailedRegReq(ctx context.Context, managerID uint64, failedReq models.FailedReq) error {
	var updateReqID uint64
	err := pr.conn.QueryRow(
		`UPDATE registration_requests
		SET (manager_id, status, message) = ($2, 'failed', $3)
		WHERE id = $1
		RETURNING id;`,
		failedReq.ReqId,
		managerID,
		failedReq.FailedMessage,
	).Scan(
		&updateReqID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (pr *postgresqlRepository) FixRegReq(ctx context.Context, reqID uint64) error {
	var updatedPassport uint64
	now := time.Now().Unix()
	err := pr.conn.QueryRow(
		`UPDATE registration_requests
		SET (status, create_time) = ('pending', $2)
		WHERE id = $1
		RETURNING id;`,
		reqID,
		now,
	).Scan(
		&updatedPassport,
	)

	if err != nil {
		return err
	}
	return nil
}
