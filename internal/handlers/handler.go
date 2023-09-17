package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lzzzzl/page-turner-pro/internal/app"
)

func RegisterHandlers(router *gin.Engine, app *app.Application) {
	registerHandlers(router, app)
}

func registerHandlers(router *gin.Engine, app *app.Application) {
	// mount all handlers under /api path
	r := router.Group("/api")
	v1 := r.Group("/v1")

	// Add health-check
	v1.GET("/health", healthCheckHandler())
}
