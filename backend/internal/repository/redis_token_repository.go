package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/go-redis/redis/v8"
)

type redisTokenRepository struct {
	Redis *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) model.TokenRepository {
	return &redisTokenRepository{
		Redis: redisClient,
	}
}

func (r *redisTokenRepository) SetRefreshToken(ctx context.Context, userID, tokenID string, expiresIn time.Duration) error {
	key := fmt.Sprintf("%s:%s", userID, tokenID)
	if err := r.Redis.Set(ctx, key, 0, expiresIn).Err(); err != nil {
		logger.Warn("could not SET refresh token to redis for userID/tokenID: %s/%s: %v", userID, tokenID, err)
		return apperrors.NewInternal()
	}
	return nil
}

func (r *redisTokenRepository) DeleteRefreshToken(ctx context.Context, userID, prevTokenID string) error {
	key := fmt.Sprintf("%s:%s", userID, prevTokenID)
	result := r.Redis.Del(ctx, key)
	if err := result.Err(); err != nil {
		logger.Warn("could not DEL refresh token to redis for userID/tokenID: %s/%s: %v", userID, prevTokenID, err)
		return apperrors.NewInternal()
	}

	if result.Val() < 1 {
		logger.Warn("refresh token to redis for userID/tokenID: %s/%s does not exists", userID, prevTokenID)
		return apperrors.NewAuthorization("invalid refresh token")
	}

	return nil
}

func (r *redisTokenRepository) DeleteUserRefreshToken(ctx context.Context, userID string) error {
	pattern := fmt.Sprintf("%s*", userID)

	iter := r.Redis.Scan(ctx, 0, pattern, 5).Iterator()
	failCount := 0

	if iter.Next(ctx) {
		if err := r.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			logger.Error("failes to delete found refrash token: %s, err: %v", iter.Val(), err)
			failCount++
		}
	}

	if err := iter.Err(); err != nil {
		logger.Warn("failes to delete refrash token: %s, err: %v", iter.Val(), err)
		failCount++
	}

	if failCount > 0 {
		return apperrors.NewInternal()
	}

	return nil

}
