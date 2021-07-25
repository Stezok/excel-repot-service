package service

import "github.com/Stezok/excel-repot-service/internal/models"

type UpdateTimeService interface {
	GetLastUpdateTime() (int64, error)
	SetLastUpdateTime(int64) error
}

type ReportService interface {
	UpdateReports() ([]models.Report, error)
	GetReports() ([]models.Report, error)
}

type Service struct {
	UpdateTimeService
	ReportService
}
