package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-project/main/app"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				app.Log.Error("Panic occurred", zap.Any("error", err))
				FailV2(c, 500, "Internal ServerConfig Error")

				c.Abort()
			}
		}()
		c.Next()
	}
}
