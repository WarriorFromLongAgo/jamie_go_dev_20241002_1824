package main

import (
	"context"
	"go-project/chain/eth"
	globalconst "go-project/common"
	"go-project/main/anvil"
	"go-project/scheduled"
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

	anvilUrl := anvil.GetAnvilURL(cfg)

	ctx := context.Background()
	ethClient, err := eth.DialEthClient(ctx, anvilUrl)
	if err != nil {
		logger.Fatal("Failed to create Ethereum client", zap.Error(err))
	}
	erc20Client, err := eth.NewTestErc20Client(ctx, anvilUrl, globalconst.TEMP_TEST_ERC20_ADDRESS)
	if err != nil {
		logger.Fatal("Failed to create NewTestErc20Client", zap.Error(err))
	}

	scanBlock, err := scheduled.NewScanBlock(ctx, ethClient, dbb, logger)
	if err != nil {
		logger.Fatal("Failed to create ScanBlock", zap.Error(err))
	}
	processingFLow, err := scheduled.NewProcessingFLow(ctx, ethClient, erc20Client, dbb, logger)
	if err != nil {
		logger.Fatal("Failed to create processingFLow", zap.Error(err))
	}
	incrementBlock, err := scheduled.NewTestIncrementBlock(ctx, ethClient, erc20Client, dbb, logger)
	if err != nil {
		logger.Fatal("Failed to create incrementBlock", zap.Error(err))
	}
	go scanBlock.Start()
	go processingFLow.Start()
	go incrementBlock.Start()

	server.RunServer(cfg, logger, dbb)
}
