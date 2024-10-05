package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	globalconst "go-project/common"
	"go-project/main/config"
	"go-project/main/log"
	"go.uber.org/zap"
	"math/big"
	"testing"
	"time"
)

func TestBusinessService_TransferERC20(t *testing.T) {
	ctx := context.Background()
	ethClient, err := DialEthClient(ctx, "http://127.0.0.1:8545")
	if err != nil {
		t.Fatalf("Failed to create EthClient: %v", err)
	}
	defer ethClient.Close()

	erc20Client, err := NewTestErc20Client(ctx, "http://127.0.0.1:8545", globalconst.TEMP_TEST_ERC20_ADDRESS)
	if err != nil {
		t.Fatalf("Failed to create TestErc20Client: %v", err)
	}

	//dsn := "root:123456@tcp(192.168.101.55:3306)/workflow_management?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	t.Fatalf("Failed to connect to database: %v", err)
	//}
	cfg := &config.Configuration{
		Log: config.LogConfig{
			Level:    "info",
			RootDir:  "./logs",
			Filename: "test.log",
			Format:   "console",
			ShowLine: true,
		},
	}
	logger, err := log.NewLogger(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	businessService := NewEthBusinessService(ethClient, erc20Client, logger)

	privateKey, err := crypto.HexToECDSA(globalconst.OWNER_PRV_KEY)
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	toAddress := common.HexToAddress(globalconst.TEMP_TO_ADDRESS)

	printERC20Balance(t, ctx, erc20Client, fromAddress, "From (before)")
	printERC20Balance(t, ctx, erc20Client, toAddress, "To (before)")

	amount := big.NewInt(9 * 1e6)
	txHash, transferData, err := businessService.TransferERC20(ctx, privateKey, fromAddress.Hex(), toAddress.Hex(), globalconst.TEMP_TEST_ERC20_ADDRESS, amount)
	if err != nil {
		t.Fatalf("Failed to transfer ERC20: %v", err)
	}
	logger.Info("txHash", zap.String("txHash", txHash))
	logger.Info("transferData", zap.String("transferData", hexutil.Encode(transferData)))

	time.Sleep(5 * time.Second)

	printERC20Balance(t, ctx, erc20Client, fromAddress, "From (after)")
	printERC20Balance(t, ctx, erc20Client, toAddress, "To (after)")

}
