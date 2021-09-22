package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Stezok/excel-repot-service/internal/models"
	"github.com/go-redis/redis/v8"
)

type RedisReportRepository struct {
	client *redis.Client
}

func (repo *RedisReportRepository) GetReportsKeys(tag string) ([]string, error) {
	ctx := context.Background()
	tag = fmt.Sprintf("report:%s", tag)
	return repo.client.HKeys(ctx, tag).Result()
}

func (repo *RedisReportRepository) GetReport(tag string, key string) (report models.Report, err error) {
	ctx := context.Background()

	tag = fmt.Sprintf("report:%s", tag)
	var jsonReport string
	jsonReport, err = repo.client.HGet(ctx, tag, key).Result()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(jsonReport), &report)
	return
}

func (repo *RedisReportRepository) SetReport(tag string, report models.Report) error {
	byteJsonReport, err := json.Marshal(report)
	if err != nil {
		return err
	}

	ctx := context.Background()
	tag = fmt.Sprintf("report:%s", tag)
	return repo.client.HSetNX(ctx, tag, report.Index, byteJsonReport).Err()
}

func (repo *RedisReportRepository) DeleteReports(tag string) error {
	ctx := context.Background()
	tag = fmt.Sprintf("report:%s", tag)
	return repo.client.Del(ctx, tag).Err()
}

func NewRedisReportRepository(client *redis.Client) *RedisReportRepository {
	return &RedisReportRepository{
		client: client,
	}
}
