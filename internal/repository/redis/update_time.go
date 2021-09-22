package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisUpdateTimeRepository struct {
	client *redis.Client
}

func (repo *RedisUpdateTimeRepository) GetLastUpdateTime(tag string) (time int64, err error) {
	ctx := context.Background()
	var result string
	tag = fmt.Sprintf("update_time:%s", tag)
	result, err = repo.client.Get(ctx, tag).Result()
	if err != nil {
		return
	}

	time, err = strconv.ParseInt(result, 10, 64)
	return
}

func (repo *RedisUpdateTimeRepository) SetLastUpdateTime(tag string, time int64) error {
	ctx := context.Background()
	tag = fmt.Sprintf("update_time:%s", tag)
	return repo.client.Set(ctx, tag, time, 0).Err()
}

func NewRedisUpdateTimeRepository(client *redis.Client) *RedisUpdateTimeRepository {
	return &RedisUpdateTimeRepository{
		client: client,
	}
}
