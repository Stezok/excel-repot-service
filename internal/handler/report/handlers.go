package report

import (
	"log"
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
		h.Logger.Print(err)
		return
	}

	ok, err := h.Service.Login(loginForm.Password)
	log.Print(loginForm.Password)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}

	if !ok {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}

	token, err := CreateToken()
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}

	ctx.SetCookie("auth_token", token, 60*60*24, "/", "", false, true)
	ctx.Redirect(http.StatusMovedPermanently, "/")
}

func (h *ReportHandler) HandlerIndex(ctx *gin.Context) {
	tag := ctx.Param("tag")
	h.Logger.Print(tag)
	intDate, err := h.Service.GetLastUpdateTime(tag)
	if err != nil {
		h.Logger.Print(err)
	}

	project := ProjectTemplate{
		ID:             tag,
		LastUpdateTime: time.Unix(intDate, 0).Format(time.ANSIC),
	}
	ctx.HTML(http.StatusOK, "index.html", project)
}

func (h *ReportHandler) HandlerData(ctx *gin.Context) {
	tag := ctx.Param("tag")
	reports, err := h.Service.GetReports(tag)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}

	var adaptedReports []ReportJSAdapter
	for _, report := range reports {
		adapter, err := NewReportJSAdapter(report)
		if err != nil {
			h.Logger.Print(err)
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
		h.Logger.Print(err)
		return
	}

	tag := ctx.Param("tag")
	filename := "temp_plan" + tag + ".xlsx"
	if err = ctx.SaveUploadedFile(uploadedFile, filename); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.Logger.Print(err)
		return
	}

	mainFilename := "plan" + tag + ".xlsx"
	err = SafeCopyFile(mainFilename, filename)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}

	err = os.Remove(filename)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}

	_, err = h.Service.UpdateReports(tag)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}
	err = h.Service.SetLastUpdateTime(tag, time.Now().Unix())
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}
}

func (h *ReportHandler) HandleUpdateReview(ctx *gin.Context) {
	uploadedFile, err := ctx.FormFile("review")
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.Logger.Print(err)
		return
	}

	tag := ctx.Param("tag")
	filename := "temp_review" + tag + ".xlsx"
	if err = ctx.SaveUploadedFile(uploadedFile, filename); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		h.Logger.Print(err)
		return
	}

	mainFilename := "review" + tag + ".xlsx"
	err = SafeCopyFile(mainFilename, filename)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}

	err = os.Remove(filename)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}

	_, err = h.Service.UpdateReports(tag)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}
	err = h.Service.SetLastUpdateTime(tag, time.Now().Unix())
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}
}

func (h *ReportHandler) HandlerLastUpdateTime(ctx *gin.Context) {
	tag := ctx.Param("tag")
	lastUpdateTime, err := h.Service.GetLastUpdateTime(tag)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		h.Logger.Print(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"last_update_time": lastUpdateTime})
}
