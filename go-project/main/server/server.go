package server

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
	"gorm.io/gorm"

	"go-project/business"
	"go-project/common/web"
	"go-project/main/config"
	"go-project/main/log"
)

func RunServer(cfg *config.Configuration, log *log.ZapLogger, db *gorm.DB) {
	ginRouter := gin.Default()

	ginRouter.Use(web.CorsHandler())
	ginRouter.Use(web.ErrorHandler(log))

	router := &business.Route{
		DB:  db,
		Log: log,
	}
	router.Register(ginRouter)

	service := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: ginRouter,
	}
	go func() {
		log.Info("Starting server", zap.String("port", cfg.Server.Port))
		if err := service.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.Shutdown(ctx); err != nil {
		log.Fatal("ServerConfig forced to shutdown", zap.Error(err))
	}

	log.Info("ServerConfig exiting")
}
