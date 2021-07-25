package report

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *ReportHandler) HandlerIndex(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func (h *ReportHandler) HandlerData(ctx *gin.Context) {
	reports, err := h.service.GetReports()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	var adaptedReports []ReportJSAdapter
	for _, report := range reports {
		adaptedReports = append(adaptedReports, NewReportJSAdapter(report))
	}

	ctx.JSON(http.StatusOK, adaptedReports)
}

func (h *ReportHandler) HandleUpdate(ctx *gin.Context) {
	reports, err := h.service.UpdateReports()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	htmlTable := h.presenter.Present(reports)
	_, err = ctx.Writer.Write([]byte(htmlTable))
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	h.service.SetLastUpdateTime(time.Now().Unix())
}

func (h *ReportHandler) HandlerLastUpdateTime(ctx *gin.Context) {
	lastUpdateTime, err := h.service.GetLastUpdateTime()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"last_update_time": lastUpdateTime})
}
