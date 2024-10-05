package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"

	globalconst "go-project/common"
	"go-project/main/log"
)

type BusinessService struct {
	ethClient   EthClient
	erc20Client TestErc20Client
	log         *log.ZapLogger
}

func NewEthBusinessService(ethClient EthClient, erc20Client TestErc20Client, log *log.ZapLogger) *BusinessService {
	return &BusinessService{
		ethClient:   ethClient,
		erc20Client: erc20Client,
		log:         log,
	}
}

func (s *BusinessService) TransferERC20(
	ctx context.Context,
	prvKey *ecdsa.PrivateKey,
	fromAddress string,
	toAddress string,
	contractAddress string,
	amount *big.Int,
) (string, []byte, error) {
	maxRetries := 3
	var lastErr error
	var txHash string
	var data []byte

	for attempt := 0; attempt < maxRetries; attempt++ {
		hash, transferData, err := s.attemptTransferERC20(ctx, prvKey, fromAddress, toAddress, contractAddress, amount)
		if err == nil {
			return hash, transferData, nil // 交易成功，返回交易哈希和data数组
		}

		lastErr = err
		s.log.Error("TransferERC20 尝试失败，准备重试", zap.Int("尝试次数", attempt+1), zap.Error(err))

		if attempt < maxRetries-1 {
			time.Sleep(3 * time.Second) // 在重试之前等待一段时间
		}
		txHash = hash       // 保存最后一次尝试的交易哈希
		data = transferData // 保存最后一次尝试的data数组
	}

	return txHash, data, fmt.Errorf("TransferERC20 在 %d 次尝试后失败: %w", maxRetries, lastErr)
}

func (s *BusinessService) attemptTransferERC20(
	ctx context.Context,
	prvKey *ecdsa.PrivateKey,
	fromAddress string,
	toAddress string,
	contractAddress string,
	amount *big.Int,
) (string, []byte, error) {
	from := common.HexToAddress(fromAddress)
	to := common.HexToAddress(toAddress)

	nonce, err := s.ethClient.TxCountByAddress(from)
	if err != nil {
		return "", nil, fmt.Errorf("获取nonce失败: %w", err)
	}

	gasPrice, err := s.ethClient.SuggestGasPrice()
	if err != nil {
		return "", nil, fmt.Errorf("获取gas价格失败: %w", err)
	}

	// 增加 gas 价格以提高交易成功率
	adjustedGasPrice := new(big.Int).Mul(gasPrice, big.NewInt(120))
	adjustedGasPrice = adjustedGasPrice.Div(adjustedGasPrice, big.NewInt(100))

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256(transferFnSignature)
	methodID := hash[:4]
	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	erc20Address := common.HexToAddress(contractAddress)
	tx := types.NewTransaction(uint64(nonce), erc20Address, big.NewInt(0), 300000, adjustedGasPrice, data)

	chainID := big.NewInt(globalconst.ChainId)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), prvKey)
	if err != nil {
		return "", data, fmt.Errorf("签名交易失败: %w", err)
	}

	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return "", data, fmt.Errorf("序列化交易失败: %w", err)
	}
	rawTxHex := hexutil.Encode(rawTxBytes)

	err = s.ethClient.SendRawTransaction(rawTxHex)
	if err != nil {
		return signedTx.Hash().Hex(), data, fmt.Errorf("发送原始交易失败: %w", err)
	}

	err = WaitForTransaction(ctx, s.ethClient, signedTx.Hash())
	if err != nil {
		return signedTx.Hash().Hex(), data, fmt.Errorf("等待交易确认失败: %w", err)
	}

	s.log.Info("TransferERC20 交易成功", zap.String("txHash", signedTx.Hash().Hex()))
	return signedTx.Hash().Hex(), data, nil
}

func WaitForTransaction(ctx context.Context, ethClient EthClient, txHash common.Hash) error {
	retries := 3
	for i := 0; i < retries; i++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("交易等待超时")
		default:
			receipt, err := ethClient.TxReceiptByTxHash(txHash)
			if err != nil {
				if err.Error() == "not found" {
					time.Sleep(2 * time.Second)
					continue
				}
				return fmt.Errorf("获取交易收据失败: %w", err)
			}
			if receipt != nil {
				if receipt.Status == types.ReceiptStatusSuccessful {
					return nil
				} else {
					return fmt.Errorf("交易失败")
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
	return fmt.Errorf("交易确认超时")
}
