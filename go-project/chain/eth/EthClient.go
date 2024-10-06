package eth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
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
	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
	BlockByNumberV2(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockByNumberV3(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockByNumberReturnJson(ctx context.Context, number *big.Int) (*types.Block, error)
	// BlockHeaderByNumber(*big.Int) (*types.Header, error)
	// LatestSafeBlockHeader() (*types.Header, error)
	LatestFinalizedBlockHeader() (*types.Header, error)
	BlockHeaderByBlockHash(common.Hash) (*types.Header, error)
	BlockHeaderListByRange(*big.Int, *big.Int) ([]*types.Header, error)

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
	rpc       rpc.RPC
	ethClient *ethclient.Client
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

	ethClient, err := ethclient.DialContext(ctx, rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create ethclient: %w", err)
	}

	return &client{
		rpc:       rpc.NewRPC(rpcClient),
		ethClient: ethClient,
	}, nil
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

func (c *client) BlockHeaderListByRange(startHeight, endHeight *big.Int) ([]*types.Header, error) {
	if startHeight.Cmp(endHeight) == 0 {
		return []*types.Header{}, nil
	}

	count := new(big.Int).Sub(endHeight, startHeight).Uint64() + 1
	headers := make([]*types.Header, count)
	batchElems := make([]gethrpc.BatchElem, count)

	for i := uint64(0); i < count; i++ {
		height := new(big.Int).Add(startHeight, new(big.Int).SetUint64(i))
		headers[i] = &types.Header{}
		batchElems[i] = gethrpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{toBlockNumArg(height), false},
			Result: headers[i],
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	err := c.rpc.BatchCallContext(ctx, batchElems)
	if err != nil {
		return nil, err
	}

	size := 0
	for _, batchElem := range batchElems {
		if batchElem.Error != nil {
			if size == 0 {
				return nil, batchElem.Error
			}
			break
		}
		size++
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

func (c *client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	var block *types.Block
	err := c.rpc.CallContext(ctx, &block, "eth_getBlockByNumber", toBlockNumArg(number), true)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, ethereum.NotFound
	}
	blockJson, _ := json.Marshal(block)
	fmt.Printf("原始 blockJson 响应: %s\n", string(blockJson))

	return block, nil
}

func (c *client) BlockByNumberV2(ctx context.Context, number *big.Int) (*types.Block, error) {
	var raw json.RawMessage
	err := c.rpc.CallContext(ctx, &raw, "eth_getBlockByNumber", toBlockNumArg(number), true)
	if err != nil {
		return nil, err
	}

	// 打印原始 JSON 响应
	fmt.Printf("原始 JSON 响应: %s\n", string(raw))

	if len(raw) == 0 {
		return nil, ethereum.NotFound
	}

	var block types.Block
	if err := json.Unmarshal(raw, &block); err != nil {
		return nil, fmt.Errorf("解析区块数据失败: %v", err)
	}

	// 使用自定义方法序列化区块
	blockJSON, err := blockToJSON(&block)
	if err != nil {
		return nil, fmt.Errorf("序列化区块数据失败: %v", err)
	}
	fmt.Printf("处理后的区块 JSON: %s\n", string(blockJSON))

	return &block, nil
}

func (c *client) BlockByNumberV3(ctx context.Context, number *big.Int) (*types.Block, error) {
	block, err := c.ethClient.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	if block == nil {
		return nil, ethereum.NotFound
	}

	blockJson, _ := json.Marshal(block)
	fmt.Printf("原始 blockJson 响应: %s\n", string(blockJson))

	return block, nil
}

// 自定义序列化方法
func blockToJSON(b *types.Block) ([]byte, error) {
	type BlockAlias types.Block
	return json.Marshal(&struct {
		*BlockAlias
		Hash common.Hash `json:"hash"`
	}{
		BlockAlias: (*BlockAlias)(b),
		Hash:       b.Hash(),
	})
}

func (c *client) BlockByNumberReturnJson(ctx context.Context, number *big.Int) (*types.Block, error) {
	var rawResponse json.RawMessage
	err := c.rpc.CallContext(ctx, &rawResponse, "eth_getBlockByNumber", toBlockNumArg(number), true)
	if err != nil {
		return nil, err
	}

	fmt.Printf("原始 RPC 响应: %s\n", string(rawResponse))

	var block *types.Block
	err = json.Unmarshal(rawResponse, &block)
	if err != nil {
		return nil, fmt.Errorf("解析区块数据失败: %v", err)
	}

	if block == nil {
		return nil, ethereum.NotFound
	}

	fmt.Printf("区块中的交易数量: %d\n", len(block.Transactions()))

	return block, nil
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
