package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"time"

	"go-project/util/retry"
	"go-project/util/rpc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

const (
	// defaultDialTimeout is default duration the processor will wait on
	// startup to make a connection to the backend
	defaultDialTimeout = 5 * time.Second

	// defaultDialAttempts is the default attempts a connection will be made
	// before failing
	defaultDialAttempts = 5

	// defaultRequestTimeout is the default duration the processor will
	// wait for a request to be fulfilled
	defaultRequestTimeout = 100 * time.Second
)

type EthClient interface {
	// BlockHeaderByNumber(*big.Int) (*types.Header, error)
	// LatestSafeBlockHeader() (*types.Header, error)
	LatestFinalizedBlockHeader() (*types.Header, error)
	BlockHeaderByBlockHash(common.Hash) (*types.Header, error)
	BlockHeaderListByRange(*big.Int, *big.Int) ([]types.Header, error)

	TxByTxHash(common.Hash) (*types.Transaction, error)

	TxReceiptByTxHash(common.Hash) (*types.Receipt, error)
	TxCountByAddress(common.Address) (hexutil.Uint64, error)
	SuggestGasPrice() (*big.Int, error)
	SuggestGasTipCap() (*big.Int, error)
	SendRawTransaction(rawTx string) error

	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)

	// Close closes the underlying RPC connection.
	// RPC close does not return any errors, but does shut down e.g. a websocket connection.
	Close()
}

type client struct {
	rpc rpc.RPC
}

func DialEthClient(ctx context.Context, rpcUrl string) (EthClient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	bOff := retry.Exponential()
	rpcClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*gethrpc.Client, error) {
		if !IsURLAvailable(rpcUrl) {
			return nil, fmt.Errorf("address unavailable (%s)", rpcUrl)
		}

		client, err := gethrpc.DialContext(ctx, rpcUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to dial address (%s): %w", rpcUrl, err)
		}
		return client, nil
	})

	if err != nil {
		return nil, err
	}

	return &client{rpc: rpc.NewRPC(rpcClient)}, nil
}

//func (c *client) BlockHeaderByNumber(number *big.Int) (*types.Header, error) {
//	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
//	defer cancel()
//
//	var header *types.Header
//	err := c.rpc.CallContext(ctxwt, &header, "eth_getBlockByNumber", toBlockNumArg(number), false)
//	if err != nil {
//		log.Error("Call eth_getBlockByNumber method fail", "err", err)
//		return nil, err
//	} else if header == nil {
//		log.Warn("header not found")
//		return nil, ethereum.NotFound
//	}
//
//	return header, nil
//}

//
//func (c *client) LatestSafeBlockHeader() (*types.Header, error) {
//	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
//	defer cancel()
//
//	var header *types.Header
//	err := c.rpc.CallContext(ctxwt, &header, "eth_getBlockByNumber", "safe", false)
//	if err != nil {
//		return nil, err
//	} else if header == nil {
//		return nil, ethereum.NotFound
//	}
//
//	return header, nil
//}

func (c *client) LatestFinalizedBlockHeader() (*types.Header, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctxwt, &header, "eth_getBlockByNumber", "finalized", false)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}

	return header, nil
}

func (c *client) BlockHeaderByBlockHash(hash common.Hash) (*types.Header, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctxwt, &header, "eth_getBlockByHash", hash, false)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}

	// sanity check on the data returned
	if header.Hash() != hash {
		return nil, errors.New("header mismatch")
	}

	return header, nil
}

func (c *client) BlockHeaderListByRange(startHeight, endHeight *big.Int) ([]types.Header, error) {
	if startHeight.Cmp(endHeight) == 0 {
		return []types.Header{}, nil
	}

	count := new(big.Int).Sub(endHeight, startHeight).Uint64() + 1
	headers := make([]types.Header, count)
	batchElems := make([]gethrpc.BatchElem, count)

	for i := uint64(0); i < count; i++ {
		height := new(big.Int).Add(startHeight, new(big.Int).SetUint64(i))
		batchElems[i] = gethrpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{toBlockNumArg(height), false},
			Result: &headers[i],
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	err := c.rpc.BatchCallContext(ctx, batchElems)
	if err != nil {
		return nil, err
	}

	size := 0
	for i, batchElem := range batchElems {
		if batchElem.Error != nil {
			if size == 0 {
				return nil, batchElem.Error
			}
			break
		}
		header, ok := batchElem.Result.(*types.Header)
		if !ok {
			return nil, fmt.Errorf("RPC response(%v) con't transfer types.Header", batchElem.Result)
		}

		headers[i] = *header
		size++

		// if i > 0 && headers[i].ParentHash != headers[i-1].Hash() {
		//     return nil, fmt.Errorf("ParentHash %s now follow %s", headers[i].Hash(), headers[i-1].Hash())
		// }
	}

	return headers[:size], nil
}

func (c *client) TxByTxHash(hash common.Hash) (*types.Transaction, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var tx *types.Transaction
	err := c.rpc.CallContext(ctxwt, &tx, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, err
	} else if tx == nil {
		return nil, ethereum.NotFound
	}

	return tx, nil
}

func (c *client) TxReceiptByTxHash(hash common.Hash) (*types.Receipt, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var txReceipt *types.Receipt
	err := c.rpc.CallContext(ctxwt, &txReceipt, "eth_getTransactionReceipt", hash)
	if err != nil {
		return nil, err
	} else if txReceipt == nil {
		return nil, ethereum.NotFound
	}

	return txReceipt, nil
}

func (c *client) TxCountByAddress(address common.Address) (hexutil.Uint64, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var nonce hexutil.Uint64
	err := c.rpc.CallContext(ctxwt, &nonce, "eth_getTransactionCount", address, "latest")
	if err != nil {
		log.Error("Call eth_getTransactionCount method fail", "err", err)
		return 0, err
	}
	log.Info("get nonce by address success", "nonce", nonce)
	return nonce, err
}

func (c *client) SuggestGasPrice() (*big.Int, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var hex hexutil.Big
	if err := c.rpc.CallContext(ctxwt, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (c *client) SuggestGasTipCap() (*big.Int, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var hex hexutil.Big
	if err := c.rpc.CallContext(ctxwt, &hex, "eth_maxPriorityFeePerGas"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

func (c *client) SendRawTransaction(rawTx string) error {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	if err := c.rpc.CallContext(ctxwt, nil, "eth_sendRawTransaction", rawTx); err != nil {
		return err
	}
	log.Info("send tx to ethereum success")
	return nil
}

func (c *client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	var result hexutil.Big
	err := c.rpc.CallContext(ctx, &result, "eth_getBalance", account, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return (*big.Int)(&result), nil
}

func (c *client) Close() {
	c.rpc.Close()
}

func IsURLAvailable(address string) bool {
	u, err := url.Parse(address)
	if err != nil {
		return false
	}
	addr := u.Host
	if u.Port() == "" {
		switch u.Scheme {
		case "http", "ws":
			addr += ":80"
		case "https", "wss":
			addr += ":443"
		default:
			// Fail open if we can't figure out what the port should be
			return true
		}
	}
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	err = conn.Close()
	if err != nil {
		return false
	}
	return true
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	return gethrpc.BlockNumber(number.Int64()).String()
}
