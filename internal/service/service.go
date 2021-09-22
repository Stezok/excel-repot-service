package service

import "github.com/Stezok/excel-repot-service/internal/models"

type AuthService interface {
	Login(string) (bool, error)
}

type UpdateTimeService interface {
	GetLastUpdateTime(string) (int64, error)
	SetLastUpdateTime(string, int64) error
}

type ReportService interface {
	UpdateReports(string) ([]models.Report, error)
	GetReports(string) ([]models.Report, error)
}

type Service struct {
	AuthService
	UpdateTimeService
	ReportService
}
