package redis

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisUpdateTimeRepository struct {
	client *redis.Client
}

func (repo *RedisUpdateTimeRepository) GetLastUpdateTime() (time int64, err error) {
	ctx := context.Background()
	var result string
	result, err = repo.client.Get(ctx, "update_time:last").Result()
	if err != nil {
		return
	}

	time, err = strconv.ParseInt(result, 10, 64)
	return
}

func (repo *RedisUpdateTimeRepository) SetLastUpdateTime(time int64) error {
	ctx := context.Background()
	return repo.client.Set(ctx, "update_time:last", time, 0).Err()
}

func NewRedisUpdateTimeRepository(client *redis.Client) *RedisUpdateTimeRepository {
	return &RedisUpdateTimeRepository{
		client: client,
	}
}
