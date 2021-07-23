package report

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ReportHandler) HandleUpdate(ctx *gin.Context) {
	reports, err := h.service.UpdateReports()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	htmlTable := h.presenter.Present(reports)
	_, err = ctx.Writer.Write([]byte(htmlTable))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (h *ReportHandler) HandleTable(ctx *gin.Context) {
	reports, err := h.service.GetReports()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	htmlTable := h.presenter.Present(reports)
	_, err = ctx.Writer.Write([]byte(htmlTable))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
