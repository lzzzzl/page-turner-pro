package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthCheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(http.StatusOK)
	}
}
