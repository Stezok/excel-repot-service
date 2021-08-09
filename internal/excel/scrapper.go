package excel

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

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
	dateTag     = "Date of SWAP"
)

func IsToday(date time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := time.Now().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

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
	columnAnchors[dateTag] = -1

	firstRow := cells[0][0]
	for i := 0; i < len(firstRow); i++ {
		if val, ok := columnAnchors[firstRow[i]]; ok && val == -1 {
			columnAnchors[firstRow[i]] = i
		}
	}

	for i := 1; i < len(cells[0]); i++ {
		row := cells[0][i]
		duid := row[columnAnchors[duidTag]]
		if _, ok := reports[duid]; !ok || duid == "" {
			continue
		}

		// date, err := time.Parse("01.02.2006", row[columnAnchors[dateTag]])
		// if err != nil || !IsToday(date) {
		// 	log.Print(err)
		// 	continue
		// }

		i := 0
		for {
			var id string
			log.Print(i)
			if i == 0 {
				id = duid
			} else {
				id = fmt.Sprintf("%d::%s", i, duid)
			}

			if val, ok := reports[id]; ok {
				if val.SQCCheck == "no owner" {
					val.SQCCheck = row[columnAnchors[sqcCheckTag]]
					reports[id] = val
				}
			} else {
				break
			}

			i++
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
			collected = 0
		}
		passed, err := strconv.Atoi(row[columnAnchors[passedTag]])
		if err != nil {
			passed = 0
		}
		failed, err := strconv.Atoi(row[columnAnchors[failedTag]])
		if err != nil {
			failed = 0
		}

		index := duid

		i := 0
		for {
			if i != 0 {
				index = fmt.Sprintf("%d::%s", i, duid)
			}

			if _, ok := reports[index]; !ok {
				break
			}
			i++
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
	log.Print(reports)
	return
}

func NewScrapper(planPath, reviewPath string) *Scrapper {
	return &Scrapper{
		planPath:   planPath,
		reviewPath: reviewPath,
	}
}
