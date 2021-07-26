package report

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *ReportHandler) HandleLoginPage(ctx *gin.Context) {
	if ctx.GetBool("registered") {
		ctx.HTML(http.StatusOK, "registered.html", nil)
		return
	}

	ctx.HTML(http.StatusOK, "login.html", nil)
}

func (h *ReportHandler) HandleLogin(ctx *gin.Context) {
	var loginForm LoginForm
	err := ctx.Bind(&loginForm)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.logger.Print(err)
		return
	}

	ok, err := h.service.Login(loginForm.Password)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	if !ok {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	token, err := CreateToken()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	ctx.SetCookie("auth_token", token, 60*60*24, "/", "localhost", false, true)
	ctx.Redirect(http.StatusMovedPermanently, "/report/")
}

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
	_, err := h.service.UpdateReports()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	h.service.SetLastUpdateTime(time.Now().Unix())
	ctx.Redirect(http.StatusMovedPermanently, "")
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
