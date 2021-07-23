package main

import (
	"github.com/Stezok/excel-repot-service/internal/excel"
	"github.com/Stezok/excel-repot-service/internal/handler/report"
	"github.com/Stezok/excel-repot-service/internal/presenter/html"
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
		ReportRepository: dbredis.NewRedisRepository(client),
	}

	scrapper := excel.NewScrapper("plan.xlsx", "review.xlsx")

	service := &service.Service{
		ReportService: service.NewCashedReportService(repo, scrapper),
	}

	presenter := html.NewHTMLPresenter()
	handler := report.NewReportHandler(service, presenter)

	router := handler.InitRoutes()
	router.Run()
}
