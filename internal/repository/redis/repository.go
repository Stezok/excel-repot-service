package redis

import (
	"context"
	"encoding/json"

	"github.com/Stezok/excel-repot-service/internal/models"
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func (repo *RedisRepository) GetReportsKeys() ([]string, error) {
	ctx := context.Background()
	return repo.client.HKeys(ctx, "report:main").Result()
}

func (repo *RedisRepository) GetReport(key string) (report models.Report, err error) {
	ctx := context.Background()

	var jsonReport string
	jsonReport, err = repo.client.HGet(ctx, "report:main", key).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(jsonReport), &report)
	return
}

func (repo *RedisRepository) SetReport(report models.Report) error {
	byteJsonReport, err := json.Marshal(report)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return repo.client.HSetNX(ctx, "report:main", report.DUID, byteJsonReport).Err()
}

func (repo *RedisRepository) DeleteReports() error {
	ctx := context.Background()
	return repo.client.Del(ctx, "report:main").Err()
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}
