package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisAuthRepository struct {
	client *redis.Client
}

func (repo *RedisAuthRepository) Login(pass string) (bool, error) {
	ctx := context.Background()
	return repo.client.SIsMember(ctx, "passwords", pass).Result()
}

func NewRedisAuthRepository(client *redis.Client) *RedisAuthRepository {
	return &RedisAuthRepository{
		client: client,
	}
}
