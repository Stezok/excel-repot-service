package excel

import (
	"errors"
	"strconv"

	"github.com/Stezok/excel-repot-service/internal/models"
	"github.com/tealeg/xlsx/v3"
)

var (
	ErrorBadFile = errors.New("bad file provided")
)

type Scrapper struct {
	planPath   string
	reviewPath string
}

const (
	duidTag     = "DU ID"
	sqcCheckTag = "SQC check"
)

func (scr *Scrapper) scrapePlan(reports map[string]models.Report) error {
	file, err := xlsx.OpenFile(scr.planPath)
	if err != nil {
		return err
	}

	cells, err := file.ToSlice()
	if err != nil {
		return err
	}

	if len(cells) == 0 {
		return ErrorBadFile
	}

	columnAnchors := make(map[string]int)
	columnAnchors[duidTag] = -1
	columnAnchors[sqcCheckTag] = -1

	lastRow := cells[0][len(cells[0])-1]
	for i := 0; i < len(lastRow); i++ {
		if val, ok := columnAnchors[lastRow[i]]; ok && val == -1 {
			columnAnchors[lastRow[i]] = i
		}
	}

	for i := 0; i < len(cells[0])-1; i++ {
		row := cells[0][i]
		duid := row[columnAnchors[duidTag]]
		if _, ok := reports[duid]; ok {
			temp := reports[duid]
			temp.SQCCheck = row[columnAnchors[sqcCheckTag]]
			reports[duid] = temp
		}
	}

	return nil
}

const (
	statusTag    = "status"
	collectedTag = "Collected Item Quantity"
	passedTag    = "Passed Item Quantity"
	failedTag    = "Failed Item Quantity"
	templateTag  = "Template Name"
)

func (scr *Scrapper) scrapeReview(reports map[string]models.Report) error {
	file, err := xlsx.OpenFile(scr.reviewPath)
	if err != nil {
		return err
	}

	cells, err := file.ToSlice()
	if err != nil {
		return err
	}
	if len(cells) == 0 {
		return ErrorBadFile
	}

	columnAnchors := make(map[string]int)
	columnAnchors[duidTag] = -1
	columnAnchors[statusTag] = -1
	columnAnchors[collectedTag] = -1
	columnAnchors[passedTag] = -1
	columnAnchors[failedTag] = -1
	columnAnchors[templateTag] = -1

	firstRow := cells[0][0]
	for i := 0; i < len(firstRow); i++ {
		if val, ok := columnAnchors[cells[0][0][i]]; ok && val == -1 {
			columnAnchors[cells[0][0][i]] = i
		}
	}

	for i := 1; i < len(cells[0]); i++ {
		row := cells[0][i]
		duid := row[columnAnchors[duidTag]]

		collected, err := strconv.Atoi(row[columnAnchors[collectedTag]])
		if err != nil {
			continue
		}
		passed, err := strconv.Atoi(row[columnAnchors[passedTag]])
		if err != nil {
			continue
		}
		failed, err := strconv.Atoi(row[columnAnchors[failedTag]])
		if err != nil {
			continue
		}

		temp := models.Report{
			DUID:         duid,
			SQCCheck:     "no owner",
			Status:       row[columnAnchors[statusTag]],
			Collected:    collected,
			Passed:       passed,
			Failed:       failed,
			TemplateName: row[columnAnchors[templateTag]],
		}

		reports[duid] = temp
	}
	return nil
}

func (scr *Scrapper) ScrapeReports() (reports []models.Report, err error) {
	reportMap := make(map[string]models.Report)
	err = scr.scrapeReview(reportMap)
	if err != nil {
		return
	}
	err = scr.scrapePlan(reportMap)
	if err != nil {
		return
	}

	for _, report := range reportMap {
		reports = append(reports, report)
	}

	return
}

func NewScrapper(planPath, reviewPath string) *Scrapper {
	return &Scrapper{
		planPath:   planPath,
		reviewPath: reviewPath,
	}
}
