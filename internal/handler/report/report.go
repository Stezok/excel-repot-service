package report

import (
	"github.com/Stezok/excel-repot-service/internal/models"
	"github.com/gin-gonic/gin"
)

type ReportService interface {
	GetReports() ([]models.Report, error)
	UpdateReports() ([]models.Report, error)
}

type Presenter interface {
	Present([]models.Report) string
}

type ReportHandler struct {
	service   ReportService
	presenter Presenter
}

func (h *ReportHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/update", h.HandleUpdate)
	router.GET("/table", h.HandleTable)

	return router
}

func NewReportHandler(service ReportService, presenter Presenter) *ReportHandler {
	return &ReportHandler{
		service:   service,
		presenter: presenter,
	}
}
