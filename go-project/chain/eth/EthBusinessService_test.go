package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"go-project/business/token/do"
	global_const "go-project/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	erc20Client, err := NewTestErc20Client(ctx, "http://127.0.0.1:8545", global_const.TEMP_TEST_ERC20_ADDRESS)
	if err != nil {
		t.Fatalf("Failed to create TestErc20Client: %v", err)
	}

	dsn := "root:123456@tcp(192.168.101.55:3306)/workflow_management?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	businessService := NewEthBusinessService(ethClient, *erc20Client, db)

	privateKey, err := crypto.HexToECDSA(global_const.OWNER_PRV_KEY)
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	toAddress := common.HexToAddress(global_const.TEMP_TO_ADDRESS)

	printERC20Balance(t, ctx, erc20Client, fromAddress, "From (before)")
	printERC20Balance(t, ctx, erc20Client, toAddress, "To (before)")

	amount := big.NewInt(9 * 1e6)
	err = businessService.TransferERC20(ctx, privateKey, 1, toAddress.Hex(), amount)
	if err != nil {
		t.Fatalf("Failed to transfer ERC20: %v", err)
	}

	time.Sleep(5 * time.Second)

	printERC20Balance(t, ctx, erc20Client, fromAddress, "From (after)")
	printERC20Balance(t, ctx, erc20Client, toAddress, "To (after)")

	var transferLog do.TokenTransferLog
	result := db.First(&transferLog, "from_address = ? AND to_address = ? AND amount = ?", fromAddress.Hex(), toAddress.Hex(), amount.Uint64())
	if result.Error != nil {
		t.Fatalf("Failed to find transfer log in database: %v", result.Error)
	}

	if transferLog.Status != do.StatusSuccess {
		t.Errorf("Expected transfer status to be %s, but got %s", do.StatusSuccess, transferLog.Status)
	}
}
