package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func (s *RedisStore) PushDeployment(ctx context.Context, deploymentId string) error {
	return s.rdb.LPush(
		ctx,
		"deployments",
		deploymentId,
	).Err()
}
func (s *RedisStore) PopDeployment(ctx context.Context) (string, error) {
	result, err := s.rdb.BRPop(
		ctx,
		0,
		"deployments",
	).Result()
	if err != nil {
		return "", err
	}

	return result[1], nil
}
