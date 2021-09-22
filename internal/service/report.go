package service

import "github.com/Stezok/excel-repot-service/internal/models"

type ReportRepository interface {
	GetReportsKeys(string) ([]string, error)
	GetReport(string, string) (models.Report, error)
	SetReport(string, models.Report) error
	DeleteReports(string) error
}

type Scrapper interface {
	ScrapeReports(string) ([]models.Report, error)
}

type CashedReportService struct {
	repo     ReportRepository
	scrapper Scrapper
}

func (rs *CashedReportService) GetReports(tag string) (reports []models.Report, err error) {
	var keys []string
	keys, err = rs.repo.GetReportsKeys(tag)
	if err != nil {
		return
	}

	for _, key := range keys {
		var report models.Report
		report, err = rs.repo.GetReport(tag, key)
		if err != nil {
			return
		}
		reports = append(reports, report)
	}

	return
}

func (rs *CashedReportService) UpdateReports(tag string) (reports []models.Report, err error) {
	reports, err = rs.scrapper.ScrapeReports(tag)
	if err != nil {
		return
	}

	err = rs.repo.DeleteReports(tag)
	if err != nil {
		return
	}

	for _, report := range reports {
		err = rs.repo.SetReport(tag, report)
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
