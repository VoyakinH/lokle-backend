package repository

import (
	"context"
	"time"

	"github.com/VoyakinH/lokle_backend/config"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type IRedisRepository interface {
	CreateSession(context.Context, string, string, time.Duration) error
	DeleteSession(context.Context, string) error
	CheckSession(context.Context, string) (string, error)
	ProlongSession(context.Context, string, time.Duration) error
}

type redisRepository struct {
	client *redis.Client
	logger logrus.Logger
}

func NewRedisRepository(cfg config.RedisConfig, logger logrus.Logger) IRedisRepository {
	return &redisRepository{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
		logger: logger,
	}
}

func (rr *redisRepository) CreateSession(ctx context.Context, sessionID string, userLogin string, expCookieTime time.Duration) error {
	_, err := rr.client.SetNX(ctx, sessionID, userLogin, expCookieTime*time.Second).Result()
	return err
}

func (rr *redisRepository) DeleteSession(ctx context.Context, cookie string) error {
	rr.client.Del(ctx, cookie).Val()
	return nil
}

func (rr *redisRepository) CheckSession(ctx context.Context, cookie string) (string, error) {
	val, err := rr.client.Get(ctx, cookie).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (rr *redisRepository) ProlongSession(ctx context.Context, cookie string, expCookieTime time.Duration) error {
	rr.client.Expire(ctx, cookie, expCookieTime*time.Second)
	return nil
}
