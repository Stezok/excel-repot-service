package report

import (
	"net/http"
	"os"
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

	ctx.SetCookie("auth_token", token, 60*60*24, "/", "", false, true)
	ctx.Redirect(http.StatusMovedPermanently, "/")
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
		adapter, err := NewReportJSAdapter(report)
		if err != nil {
			h.logger.Print(err)
			continue
		}

		adaptedReports = append(adaptedReports, adapter)
	}

	ctx.JSON(http.StatusOK, adaptedReports)
}

func (h *ReportHandler) HandleUpdatePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "update.html", nil)
}

func (h *ReportHandler) HandleUpdatePlan(ctx *gin.Context) {
	uploadedFile, err := ctx.FormFile("plan")
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.logger.Print(err)
		return
	}

	filename := "temp_plan.xlsx"
	if err = ctx.SaveUploadedFile(uploadedFile, filename); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.logger.Print(err)
		return
	}

	err = SafeCopyFile("plan.xlsx", "temp_plan.xlsx")
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	err = os.Remove("temp_plan.xlsx")
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	_, err = h.service.UpdateReports()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}
	err = h.service.SetLastUpdateTime(time.Now().Unix())
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}
}

func (h *ReportHandler) HandleUpdateReview(ctx *gin.Context) {
	uploadedFile, err := ctx.FormFile("review")
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.logger.Print(err)
		return
	}

	filename := "temp_review.xlsx"
	if err = ctx.SaveUploadedFile(uploadedFile, filename); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.logger.Print(err)
		return
	}

	err = SafeCopyFile("review.xlsx", "temp_review.xlsx")
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	err = os.Remove("temp_review.xlsx")
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}

	_, err = h.service.UpdateReports()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}
	err = h.service.SetLastUpdateTime(time.Now().Unix())
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Print(err)
		return
	}
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
