package main

import (
	"context"
	"go-project/chain/eth"
	"go-project/synchronizer"
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

	ctx := context.Background()
	ethClient, err := eth.DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		logger.Fatal("Failed to create Ethereum client", zap.Error(err))
	}

	sync, err := synchronizer.NewSynchronizer(ctx, ethClient, dbb)
	if err != nil {
		logger.Fatal("Failed to create Synchronizer", zap.Error(err))
	}
	go sync.Start()

	server.RunServer(cfg, logger, dbb)
}
