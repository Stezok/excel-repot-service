package repository

import "github.com/Stezok/excel-repot-service/internal/models"

type AuthRepository interface {
	Login(string) (bool, error)
}

type UpdateTimeRepository interface {
	GetLastUpdateTime(string) (int64, error)
	SetLastUpdateTime(string, int64) error
}

type ReportRepository interface {
	GetReportsKeys(string) ([]string, error)
	GetReport(string, string) (models.Report, error)
	SetReport(string, models.Report) error
	DeleteReports(string) error
}

type Repository struct {
	AuthRepository
	UpdateTimeRepository
	ReportRepository
}
