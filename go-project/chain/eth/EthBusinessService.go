package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"gorm.io/gorm"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"

	"go-project/business/token/do"
	globalconst "go-project/common"
)

type BusinessService struct {
	ethClient   EthClient
	erc20Client TestErc20Client
	db          *gorm.DB
}

func NewEthBusinessService(ethClient EthClient, erc20Client TestErc20Client, db *gorm.DB) *BusinessService {
	return &BusinessService{
		ethClient:   ethClient,
		erc20Client: erc20Client,
		db:          db,
	}
}

func (s *BusinessService) TransferERC20(
	ctx context.Context,
	prvKey *ecdsa.PrivateKey,
	tokenInfoID int,
	toAddress string,
	amount *big.Int,
) error {
	fromAddress := crypto.PubkeyToAddress(prvKey.PublicKey)
	to := common.HexToAddress(toAddress)

	nonce, err := s.ethClient.TxCountByAddress(fromAddress)
	if err != nil {
		return fmt.Errorf("TransferERC20 get nonce: %w", err)
	}

	gasPrice, err := s.ethClient.SuggestGasPrice()
	if err != nil {
		return fmt.Errorf("TransferERC20 get gas price: %w", err)
	}

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.Keccak256(transferFnSignature)
	methodID := hash[:4]
	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	erc20Address := common.HexToAddress(globalconst.TEMP_TEST_ERC20_ADDRESS)
	tx := types.NewTransaction(uint64(nonce), erc20Address, big.NewInt(0), 300000, gasPrice, data)

	chainID := big.NewInt(globalconst.ChainId)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), prvKey)
	if err != nil {
		return fmt.Errorf("TransferERC20 sign transaction: %w", err)
	}

	rawTxBytes, err := signedTx.MarshalBinary()
	if err != nil {
		return fmt.Errorf("TransferERC20 serialize transaction: %w", err)
	}
	rawTxHex := hexutil.Encode(rawTxBytes)

	err = s.ethClient.SendRawTransaction(rawTxHex)
	if err != nil {
		return fmt.Errorf("TransferERC20 send raw transaction: %w", err)
	}

	transferLog := &do.TokenTransferLog{
		TokenInfoID:     tokenInfoID,
		FromAddress:     fromAddress.Hex(),
		ToAddress:       toAddress,
		Amount:          amount.Uint64(),
		Status:          do.StatusPending,
		TransactionHash: signedTx.Hash().Hex(),
		CreateBy:        globalconst.SystemUser,
		CreateAddr:      fromAddress.Hex(),
		CreatedTime:     time.Now(),
	}
	tokenTransferLogManager := do.NewTokenTransferLogManager(s.db)
	err = tokenTransferLogManager.Create(transferLog)
	if err != nil {
		return fmt.Errorf("create TokenTransferLog: %w", err)
	}

	err = WaitForTransaction(ctx, s.ethClient, signedTx.Hash())
	if err != nil {
		log.Error("TransferERC20 Transfer waitForTransaction", "err", err)
		err := tokenTransferLogManager.UpdateStatus(transferLog.ID, do.StatusFailed, "", "")
		if err != nil {
			return fmt.Errorf("TransferERC20 Transfer waitForTransaction: %w", err)
		}
		return nil
	}

	err = tokenTransferLogManager.UpdateStatus(transferLog.ID, do.StatusSuccess, "", "")
	if err != nil {
		return fmt.Errorf("TransferERC20 Transfer waitForTransaction: %w", err)
	}
	log.Info("TransferERC20 Transfer success", "txHash", signedTx.Hash().Hex())
	return nil
}

func WaitForTransaction(ctx context.Context, ethClient EthClient, txHash common.Hash) error {
	retries := 30
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
