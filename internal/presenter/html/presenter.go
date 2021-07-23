package html

import (
	"fmt"
	"strconv"

	"github.com/Stezok/excel-repot-service/internal/models"
)

type HTMLPresenter struct {
}

func (pres *HTMLPresenter) Present(reports []models.Report) string {

	container := `
	<style>
		.styled-table {
			border-collapse: collapse;
			margin: 25px 0;
			font-size: 0.9em;
			font-family: sans-serif;
			min-width: 400px;
			box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
		}

		.styled-table thead tr {
			background-color: #009879;
			color: #ffffff;
			text-align: left;
		}

		.styled-table th,
		.styled-table td {
			padding: 12px 15px;
		}

		.styled-table tbody tr {
			border-bottom: 1px solid #dddddd;
		}
		
		.styled-table tbody tr:nth-of-type(even) {
			background-color: #f3f3f3;
		}
		
		.styled-table tbody tr:last-of-type {
			border-bottom: 2px solid #009879;
		}

		.styled-table tbody tr.active-row {
			font-weight: bold;
			color: #009879;
		}
	</style>

	<table class="styled-table">
		<thead>
			<tr>
				<th>DU ID</th>
				<th>SQC Check</th>
				<th>Status</th>
				<th>Template Name</th>
				<th>Collected</th>
				<th>Passed</th>
				<th>Failed</th>
				<th>Unchecked</th>
			</tr>
		</thead>
		<tbody>
			%s
		</tbody>
	</table>
	`

	body := ``
	rowPattern := `
		<tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%d</td>
			<td>%d</td>
			<td>%d</td>
			<td>%s</td>
		</tr>
	`

	for _, report := range reports {
		var unchecked string
		if report.Passed+report.Failed == report.Collected {
			unchecked = "+1"
		} else {
			unchecked = strconv.Itoa(report.Collected - report.Failed - report.Passed)
		}

		body += fmt.Sprintf(rowPattern,
			report.DUID,
			report.SQCCheck,
			report.Status,
			report.TemplateName,
			report.Collected,
			report.Passed,
			report.Failed,
			unchecked,
		)
	}

	return fmt.Sprintf(container, body)
}

func NewHTMLPresenter() *HTMLPresenter {
	return &HTMLPresenter{}
}
