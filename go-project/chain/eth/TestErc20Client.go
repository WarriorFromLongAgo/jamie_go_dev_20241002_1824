package eth

import (
	"context"
	"go-project/abigo"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TestErc20Client interface {
	BalanceOf(ctx context.Context, address common.Address) (*big.Int, error)
	Approve(auth *bind.TransactOpts, spender common.Address, amount *big.Int) (common.Hash, error)
	Transfer(auth *bind.TransactOpts, to common.Address, amount *big.Int) (common.Hash, error)
	Close() error
}

type erc20Client struct {
	instance  *abigo.Testerc20
	ethClient *ethclient.Client
}

func NewTestErc20Client(ctx context.Context, rpcUrl string, contractAddress string) (TestErc20Client, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	ethClient, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, err
	}

	address := common.HexToAddress(contractAddress)
	instance, err := abigo.NewTesterc20(address, ethClient)
	if err != nil {
		ethClient.Close()
		return nil, err
	}

	return &erc20Client{
		instance:  instance,
		ethClient: ethClient,
	}, nil
}

func (c *erc20Client) BalanceOf(ctx context.Context, address common.Address) (*big.Int, error) {
	return c.instance.BalanceOf(&bind.CallOpts{Context: ctx}, address)
}

func (c *erc20Client) Approve(auth *bind.TransactOpts, spender common.Address, amount *big.Int) (common.Hash, error) {
	tx, err := c.instance.Approve(auth, spender, amount)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (c *erc20Client) Transfer(auth *bind.TransactOpts, to common.Address, amount *big.Int) (common.Hash, error) {
	tx, err := c.instance.Transfer(auth, to, amount)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (c *erc20Client) Close() error {
	if c.ethClient != nil {
		c.ethClient.Close()
	}
	return nil
}
