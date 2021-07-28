package report

import (
	"time"

	"github.com/Stezok/excel-repot-service/internal/models"
)

type ReportJSAdapter struct {
	DUID         string `json:"du_id"`
	SQCCheck     string `json:"sqc_check"`
	Status       string `json:"status"`
	TemplateName string `json:"template_name"`
	Unchecked    int    `json:"unchecked"`
	SaveTime     int64  `json:"save_time"`
}

const layout = "2006-01-02 15:04:05"

func NewReportJSAdapter(report models.Report) (adapter ReportJSAdapter, err error) {
	saveStartTime, err := time.Parse(layout, report.SaveStartTime)
	if err != nil {
		return
	}

	saveEndTime, _ := time.Parse(layout, report.SaveEndTime)
	if err != nil {
		return
	}

	adapter = ReportJSAdapter{
		DUID:         report.DUID,
		SQCCheck:     report.SQCCheck,
		Status:       report.Status,
		TemplateName: report.TemplateName,
		Unchecked:    report.Collected - (report.Passed + report.Failed),
		SaveTime:     saveEndTime.Unix() - saveStartTime.Unix(),
	}
	return
}
