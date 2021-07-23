package presenter

import "github.com/Stezok/excel-repot-service/internal/models"

type Presenter interface {
	Present([]models.Report) string
}
