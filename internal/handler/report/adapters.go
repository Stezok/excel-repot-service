package report

import "github.com/Stezok/excel-repot-service/internal/models"

type ReportJSAdapter struct {
	DUID         string `json:"du_id"`
	SQCCheck     string `json:"sqc_check"`
	Status       string `json:"status"`
	TemplateName string `json:"template_name"`
	Unchecked    int    `json:"unchecked"`
}

func NewReportJSAdapter(report models.Report) ReportJSAdapter {
	return ReportJSAdapter{
		DUID:         report.DUID,
		SQCCheck:     report.SQCCheck,
		Status:       report.Status,
		TemplateName: report.TemplateName,
		Unchecked:    report.Collected - (report.Passed + report.Failed),
	}
}
