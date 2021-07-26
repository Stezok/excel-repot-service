package main

import (
	"log"

	"github.com/Stezok/excel-repot-service/internal/excel"
	"github.com/Stezok/excel-repot-service/internal/handler/report"
	"github.com/Stezok/excel-repot-service/internal/repository"
	dbredis "github.com/Stezok/excel-repot-service/internal/repository/redis"
	"github.com/Stezok/excel-repot-service/internal/service"
	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})

	repo := repository.Repository{
		AuthRepository:       dbredis.NewRedisAuthRepository(client),
		UpdateTimeRepository: dbredis.NewRedisUpdateTimeRepository(client),
		ReportRepository:     dbredis.NewRedisReportRepository(client),
	}

	scrapper := excel.NewScrapper("plan.xlsx", "review.xlsx")

	service := &service.Service{
		AuthService:       service.NewDefaultAuthService(repo),
		UpdateTimeService: service.NewDefaultUpdateTimeService(repo),
		ReportService:     service.NewCashedReportService(repo, scrapper),
	}

	handler := report.NewReportHandler(log.Default(), service)

	router := handler.InitRoutes()
	router.Run()
}
