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

	firstRow := cells[0][0]
	for i := 0; i < len(firstRow); i++ {
		if val, ok := columnAnchors[firstRow[i]]; ok && val == -1 {
			columnAnchors[firstRow[i]] = i
		}
	}

	for i := 1; i < len(cells[0]); i++ {
		row := cells[0][i]
		duid := row[columnAnchors[duidTag]]
		if duid == "" {
			continue
		}

		if _, ok := reports[duid]; ok {
			temp := reports[duid]

			if temp.SQCCheck != "no owner" {
				if _, ok := reports["second::"+duid]; !ok {
					temp.Index = "second::" + temp.Index
					temp.SQCCheck = row[columnAnchors[sqcCheckTag]]
					reports[temp.Index] = temp
					continue
				} else {
					temp = reports["second::"+duid]
					temp.SQCCheck = row[columnAnchors[sqcCheckTag]]
					reports["second::"+duid] = temp
				}
			} else {
				temp.SQCCheck = row[columnAnchors[sqcCheckTag]]
				reports[duid] = temp
			}
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
	saveStartTag = "Save Start Time"
	saveEndTag   = "Save End Time"
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
	columnAnchors[saveStartTag] = -1
	columnAnchors[saveEndTag] = -1

	firstRow := cells[0][0]
	for i := 0; i < len(firstRow); i++ {
		if val, ok := columnAnchors[cells[0][0][i]]; ok && val == -1 {
			columnAnchors[cells[0][0][i]] = i
		}
	}

	for i := 1; i < len(cells[0]); i++ {
		row := cells[0][i]
		duid := row[columnAnchors[duidTag]]
		if duid == "" {
			continue
		}

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

		index := duid
		if _, ok := reports[index]; ok {
			index = "second:" + duid
		}

		temp := models.Report{
			Index:         index,
			DUID:          duid,
			SQCCheck:      "no owner",
			Status:        row[columnAnchors[statusTag]],
			Collected:     collected,
			Passed:        passed,
			Failed:        failed,
			TemplateName:  row[columnAnchors[templateTag]],
			SaveStartTime: row[columnAnchors[saveStartTag]],
			SaveEndTime:   row[columnAnchors[saveEndTag]],
		}

		reports[index] = temp
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
