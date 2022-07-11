package repository

import (
	"context"
	"time"

	"github.com/VoyakinH/lokle_backend/config"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type IRedisSessionRepository interface {
	CreateSession(context.Context, string, string, time.Duration) error
	DeleteSession(context.Context, string) error
	CheckSession(context.Context, string) (string, error)
	ProlongSession(context.Context, string, time.Duration) error
}

type redisSessionRepository struct {
	client *redis.Client
	logger logrus.Logger
}

func NewRedisSessionRepository(cfg config.RedisConfig, logger logrus.Logger) IRedisSessionRepository {
	return &redisSessionRepository{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
		logger: logger,
	}
}

func (rsr *redisSessionRepository) CreateSession(ctx context.Context, sessionID string, email string, expCookieTime time.Duration) error {
	_, err := rsr.client.SetNX(ctx, sessionID, email, expCookieTime*time.Second).Result()
	return err
}

func (rsr *redisSessionRepository) DeleteSession(ctx context.Context, cookie string) error {
	rsr.client.Del(ctx, cookie).Val()
	return nil
}

func (rsr *redisSessionRepository) CheckSession(ctx context.Context, cookie string) (string, error) {
	val, err := rsr.client.Get(ctx, cookie).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (rsr *redisSessionRepository) ProlongSession(ctx context.Context, cookie string, expCookieTime time.Duration) error {
	rsr.client.Expire(ctx, cookie, expCookieTime*time.Second)
	return nil
}
