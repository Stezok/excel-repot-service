package redis

import (
	"context"
	"encoding/json"

	"github.com/Stezok/excel-repot-service/internal/models"
	"github.com/go-redis/redis/v8"
)

type RedisReportRepository struct {
	client *redis.Client
}

func (repo *RedisReportRepository) GetReportsKeys() ([]string, error) {
	ctx := context.Background()
	return repo.client.HKeys(ctx, "report:main").Result()
}

func (repo *RedisReportRepository) GetReport(key string) (report models.Report, err error) {
	ctx := context.Background()

	var jsonReport string
	jsonReport, err = repo.client.HGet(ctx, "report:main", key).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(jsonReport), &report)
	return
}

func (repo *RedisReportRepository) SetReport(report models.Report) error {
	byteJsonReport, err := json.Marshal(report)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return repo.client.HSetNX(ctx, "report:main", report.DUID, byteJsonReport).Err()
}

func (repo *RedisReportRepository) DeleteReports() error {
	ctx := context.Background()
	return repo.client.Del(ctx, "report:main").Err()
}

func NewRedisReportRepository(client *redis.Client) *RedisReportRepository {
	return &RedisReportRepository{
		client: client,
	}
}
