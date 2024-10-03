package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-project/business"
	"go-project/common/web"
)

func RunServer(cfg *Configuration) {
	ginRouter := gin.Default()

	ginRouter.Use(web.CorsHandler())
	ginRouter.Use(web.ErrorHandler())

	apiGroup := ginRouter.Group("")
	business.SetApiGroupRoutes(apiGroup)

	service := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: ginRouter,
	}
	go func() {
		Log.Info("Starting server", zap.String("port", cfg.Server.Port))
		if err := service.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			Log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.Shutdown(ctx); err != nil {
		Log.Fatal("ServerConfig forced to shutdown", zap.Error(err))
	}

	Log.Info("ServerConfig exiting")
}
