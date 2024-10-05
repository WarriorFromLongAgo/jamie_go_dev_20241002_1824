package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	global_const "go-project/common"
	"math/big"
	"testing"
	"time"
)

func TestTestErc20Client_BalanceOf(t *testing.T) {
	ctx := context.Background()
	erc20Client, err := NewTestErc20Client(ctx, "http://127.0.0.1:8545", "0x700b6A60ce7EaaEA56F065753d8dcB9653dbAD35")
	if err != nil {
		t.Fatalf("Failed to create TestErc20Client: %v", err)
	}
	address := common.HexToAddress("0xa0Ee7A142d267C1f36714E4a8F75612F20a79720")
	balanceOf, err := erc20Client.BalanceOf(ctx, address)
	if err != nil {
		t.Fatalf("Failed to get balance: %v", err)
	}
	t.Logf("TestTestErc20Client_BalanceOf: %s", balanceOf.String())
}

func TestTestErc20Client_ApproveAndTransfer(t *testing.T) {
	ctx := context.Background()
	erc20Client, err := NewTestErc20Client(ctx, "http://127.0.0.1:8545", "0x700b6A60ce7EaaEA56F065753d8dcB9653dbAD35")
	if err != nil {
		t.Fatalf("Failed to create TestErc20Client: %v", err)
	}

	fromAddress := common.HexToAddress("0xa0Ee7A142d267C1f36714E4a8F75612F20a79720")
	toAddress := common.HexToAddress("0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f")

	printERC20Balance(t, ctx, erc20Client, fromAddress, "From")
	printERC20Balance(t, ctx, erc20Client, toAddress, "To")

	privateKey, err := crypto.HexToECDSA(global_const.OWNER_PRV_KEY)
	if err != nil {
		t.Fatalf("Failed to create private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(global_const.ChainId))
	if err != nil {
		t.Fatalf("Failed to create auth: %v", err)
	}

	amount := big.NewInt(9 * 1e6)
	approveHash, err := erc20Client.Approve(auth, toAddress, amount)
	if err != nil {
		t.Fatalf("Failed to approve: %v", err)
	}
	t.Logf("Approve transaction hash: %s", approveHash.Hex())

	time.Sleep(5 * time.Second)

	transferHash, err := erc20Client.Transfer(auth, toAddress, amount)
	if err != nil {
		t.Fatalf("Failed to transfer: %v", err)
	}
	t.Logf("Transfer transaction hash: %s", transferHash.Hex())

	time.Sleep(5 * time.Second)

	printERC20Balance(t, ctx, erc20Client, fromAddress, "From")
	printERC20Balance(t, ctx, erc20Client, toAddress, "To")
}

func printERC20Balance(t *testing.T, ctx context.Context, client TestErc20Client, address common.Address, label string) {
	balance, err := client.BalanceOf(ctx, address)
	if err != nil {
		t.Fatalf("Failed to get balance for %s address: %v", label, err)
	}
	t.Logf("%s address balance: %s", label, balance.String())
}
