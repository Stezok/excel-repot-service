package report

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *ReportHandler) MiddlewareAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("auth_token")
		if err != nil {
			ctx.Next()
			return
		}

		if ok, err := IsValidToken(token); !ok || err != nil {
			h.Logger.Print(err)
			ctx.Next()
			return
		}

		ctx.Set("registered", true)
		ctx.Next()
	}
}

func (h *ReportHandler) MiddlewareAuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		registered := ctx.GetBool("registered")

		if !registered {
			ctx.Redirect(http.StatusMovedPermanently, "/login")
		}
	}
}
