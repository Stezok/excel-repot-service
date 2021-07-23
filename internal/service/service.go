package service

import "github.com/Stezok/excel-repot-service/internal/models"

type ReportService interface {
	UpdateReports() ([]models.Report, error)
	GetReports() ([]models.Report, error)
}

type Service struct {
	ReportService
}
