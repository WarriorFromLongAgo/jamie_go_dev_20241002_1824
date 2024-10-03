package eth

import (
	"context"
	"math/big"

	"go-project/abigo"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TestErc20Client struct {
	instance *abigo.Testerc20
}

func NewTestErc20Client(ctx context.Context, rpcUrl string, contractAddress string) (*TestErc20Client, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(contractAddress)
	instance, err := abigo.NewTesterc20(address, client)
	if err != nil {
		return nil, err
	}

	return &TestErc20Client{
		instance: instance,
	}, nil
}

func (c *TestErc20Client) BalanceOf(ctx context.Context, address common.Address) (*big.Int, error) {
	return c.instance.BalanceOf(&bind.CallOpts{Context: ctx}, address)
}

func (c *TestErc20Client) Approve(auth *bind.TransactOpts, spender common.Address, amount *big.Int) (common.Hash, error) {
	tx, err := c.instance.Approve(auth, spender, amount)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (c *TestErc20Client) Transfer(auth *bind.TransactOpts, to common.Address, amount *big.Int) (common.Hash, error) {
	tx, err := c.instance.Transfer(auth, to, amount)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}
