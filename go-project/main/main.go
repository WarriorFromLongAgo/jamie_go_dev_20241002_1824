package main

import (
	"go.uber.org/zap"

	"go-project/main/config"
	"go-project/main/db"
	"go-project/main/log"
	"go-project/main/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger, err := log.NewLogger(cfg)
	if err != nil {
		panic(err)
	}
	logger.Info("NewLogger success")

	dbb := db.InitializeDB(cfg, logger)
	logger.Info("InitializeDB success")

	defer func() {
		if dbb != nil {
			tempDb, _ := dbb.DB()
			err = tempDb.Close()
			if err != nil {
				logger.Error("main db.Close", zap.Any("err", err))
				return
			}
			logger.Info("main db.Close success")
		}
	}()

	server.RunServer(cfg, logger, dbb)
}
