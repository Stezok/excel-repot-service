package report

import (
	"github.com/Stezok/excel-repot-service/internal/models"
	"github.com/gin-gonic/gin"
)

type Logger interface {
	Print(...interface{})
}

type UpdateTimeService interface {
	GetLastUpdateTime() (int64, error)
	SetLastUpdateTime(int64) error
}

type ReportService interface {
	GetReports() ([]models.Report, error)
	UpdateReports() ([]models.Report, error)
}

type Service interface {
	UpdateTimeService
	ReportService
}

type Presenter interface {
	Present([]models.Report) string
}

type ReportHandler struct {
	logger    Logger
	service   Service
	presenter Presenter
}

func (h *ReportHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.LoadHTMLGlob("../../assets/*.html")

	router.GET("/", h.HandlerIndex)
	router.GET("/data", h.HandlerData)
	router.GET("/update", h.HandleUpdate)
	router.GET("/lastUpdateTime", h.HandlerLastUpdateTime)

	return router
}

func NewReportHandler(logger Logger, service Service, presenter Presenter) *ReportHandler {
	return &ReportHandler{
		logger:    logger,
		service:   service,
		presenter: presenter,
	}
}
