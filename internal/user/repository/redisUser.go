package repository

import (
	"context"
	"time"

	"github.com/VoyakinH/lokle_backend/config"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type IRedisUserRepository interface {
	AddUserToken(context.Context, string, string, time.Duration) error
	GetUserAndDelete(context.Context, string) (string, error)
}

type redisUserRepository struct {
	client *redis.Client
	logger logrus.Logger
}

func NewRedisUserRepository(cfg config.RedisConfig, logger logrus.Logger) IRedisUserRepository {
	return &redisSessionRepository{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Password: cfg.Password,
			DB:       cfg.DB,
		}),
		logger: logger,
	}
}

func (rur *redisSessionRepository) AddUserToken(ctx context.Context, token string, email string, expTokenTime time.Duration) error {
	_, err := rur.client.SetNX(ctx, token, email, expTokenTime*time.Second).Result()
	return err
}

func (rur *redisSessionRepository) GetUserAndDelete(ctx context.Context, cookie string) (string, error) {
	val, err := rur.client.Get(ctx, cookie).Result()
	if err != nil {
		return "", err
	}
	rur.client.Del(ctx, cookie).Val()
	return val, nil
}
