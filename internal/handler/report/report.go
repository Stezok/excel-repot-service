package report

import (
	"github.com/Stezok/excel-repot-service/internal/models"
	"github.com/gin-gonic/gin"
)

type Logger interface {
	Print(...interface{})
}

type AuthService interface {
	Login(string) (bool, error)
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
	AuthService
	UpdateTimeService
	ReportService
}

type ReportHandler struct {
	logger  Logger
	service Service
}

func (h *ReportHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Static("/static", "../../assets")
	router.LoadHTMLGlob("../../assets/html/*.html")

	router.Use(h.MiddlewareAuth())
	router.GET("/login", h.HandleLoginPage)
	router.POST("/login", h.HandleLogin)

	auth := router.Group("/report", h.MiddlewareAuthRequired())
	{
		auth.GET("/", h.HandlerIndex)
		auth.GET("/data", h.HandlerData)
		auth.GET("/update", h.HandleUpdatePage)
		auth.POST("/update/plan", h.HandleUpdatePlan)
		auth.POST("/update/review", h.HandleUpdateReview)
		auth.GET("/lastUpdateTime", h.HandlerLastUpdateTime)
	}

	return router
}

func NewReportHandler(logger Logger, service Service) *ReportHandler {
	return &ReportHandler{
		logger:  logger,
		service: service,
	}
}
