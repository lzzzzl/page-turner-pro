package middleware

import (
	"context"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

// SetGeneralMiddlewares add general-purpose middlewares
func SetGeneralMiddlewares(ctx context.Context, ginRouter *gin.Engine) {
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(requestid.New())
}
