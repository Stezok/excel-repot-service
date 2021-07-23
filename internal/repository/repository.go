package repository

import "github.com/Stezok/excel-repot-service/internal/models"

type ReportRepository interface {
	GetReportsKeys() ([]string, error)
	GetReport(string) (models.Report, error)
	SetReport(models.Report) error
	DeleteReports() error
}

type Repository struct {
	ReportRepository
}
