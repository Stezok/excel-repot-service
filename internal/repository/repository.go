package repository

import "github.com/Stezok/excel-repot-service/internal/models"

type AuthRepository interface {
	Login(string) (bool, error)
}

type UpdateTimeRepository interface {
	GetLastUpdateTime() (int64, error)
	SetLastUpdateTime(int64) error
}

type ReportRepository interface {
	GetReportsKeys() ([]string, error)
	GetReport(string) (models.Report, error)
	SetReport(models.Report) error
	DeleteReports() error
}

type Repository struct {
	AuthRepository
	UpdateTimeRepository
	ReportRepository
}
