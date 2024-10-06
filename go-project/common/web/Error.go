package web

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-project/main/log"
)

func ErrorHandler(logger *log.ZapLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				logger.Error("Panic occurred", zap.Any("error", err))
				FailV2(c, 500, "Internal ServerConfig Error")

				c.Abort()
			}
		}()
		c.Next()
	}
}
