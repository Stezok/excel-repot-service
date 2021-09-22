package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Stezok/excel-repot-service/internal/automatic/updater"
	"github.com/Stezok/excel-repot-service/internal/config"
	"github.com/Stezok/excel-repot-service/internal/excel"
	"github.com/Stezok/excel-repot-service/internal/handler/report"
	"github.com/Stezok/excel-repot-service/internal/repository"
	dbredis "github.com/Stezok/excel-repot-service/internal/repository/redis"
	"github.com/Stezok/excel-repot-service/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var conf config.Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal(err)
	}
	conf.PushToOSEnv()

	redisAddr := fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port)
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	status := client.Ping(context.Background())
	log.Print(status)

	repo := repository.Repository{
		AuthRepository:       dbredis.NewRedisAuthRepository(client),
		UpdateTimeRepository: dbredis.NewRedisUpdateTimeRepository(client),
		ReportRepository:     dbredis.NewRedisReportRepository(client),
	}

	scrapper := excel.NewScrapper(conf.App.PlanPath, conf.App.ReviewPath)

	service := &service.Service{
		AuthService:       service.NewDefaultAuthService(repo),
		UpdateTimeService: service.NewDefaultUpdateTimeService(repo),
		ReportService:     service.NewCashedReportService(repo, scrapper),
	}

	updaterConf := updater.UpdaterConfig{
		SeleniumPath: conf.Updater.SeleniumPath,
		Mode:         conf.Updater.BrowserMode,
		Port:         conf.Updater.Port,
		ReviewPath:   conf.App.ReviewPath,
		DownloadPath: conf.Updater.DownloadPath,
	}
	upd := updater.NewUpdater(updaterConf, service.ReportService, service.UpdateTimeService)
	go upd.Run()

	go func() {
		delayBetween := time.Duration(10 * int(time.Minute) / len(conf.Updater.HuaweiAccounts))
		for _, project := range conf.Updater.HuaweiAccounts {
			job := updater.UpdateJob{
				HuaweiLogin:    project.HuaweiLogin,
				HuaweiPassword: project.HuaweiPassword,
				ProjectID:      project.ProjectID,
			}
			upd.PushJobEvery(job, 10*time.Minute)
			time.Sleep(delayBetween)
		}
	}()

	handler := &report.ReportHandler{
		Logger:  log.Default(),
		Service: service,
		ReportHandlerConfig: report.ReportHandlerConfig{
			PathToStatic:   conf.Server.PathToStatic,
			PathToHTMLGlob: conf.Server.PathToHTMLGlob,
		},
	}
	addr := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	handler.InitRoutes().Run(addr)
}
