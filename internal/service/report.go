package service

import "github.com/Stezok/excel-repot-service/internal/models"

type ReportRepository interface {
	GetReportsKeys() ([]string, error)
	GetReport(string) (models.Report, error)
	SetReport(models.Report) error
	DeleteReports() error
}

type Scrapper interface {
	ScrapeReports() ([]models.Report, error)
}

type CashedReportService struct {
	repo     ReportRepository
	scrapper Scrapper
}

func (rs *CashedReportService) GetReports() (reports []models.Report, err error) {
	var keys []string
	keys, err = rs.repo.GetReportsKeys()
	if err != nil {
		return
	}

	for _, key := range keys {
		var report models.Report
		report, err = rs.repo.GetReport(key)
		if err != nil {
			return
		}
		reports = append(reports, report)
	}

	return
}

func (rs *CashedReportService) UpdateReports() (reports []models.Report, err error) {
	reports, err = rs.scrapper.ScrapeReports()
	if err != nil {
		return
	}

	err = rs.repo.DeleteReports()
	if err != nil {
		return
	}

	for _, report := range reports {
		err = rs.repo.SetReport(report)
		if err != nil {
			return
		}
	}

	return
}

func NewCashedReportService(repo ReportRepository, scrapper Scrapper) *CashedReportService {
	return &CashedReportService{
		repo:     repo,
		scrapper: scrapper,
	}
}
