package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
}

func NewEthBusinessService(ethClient EthClient, erc20Client TestErc20Client) *BusinessService {
	return &BusinessService{
		ethClient:   ethClient,
		erc20Client: erc20Client,
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

	auth, err := bind.NewKeyedTransactorWithChainID(prvKey, big.NewInt(globalconst.ChainId))
	if err != nil {
		return fmt.Errorf("TransferERC20 NewKeyedTransactorWithChainID: %w", err)
	}

	approveHash, err := s.erc20Client.Approve(auth, to, amount)
	if err != nil {
		return fmt.Errorf("TransferERC20 Approve: %w", err)
	}
	err = s.waitForTransaction(ctx, approveHash)
	if err != nil {
		return fmt.Errorf("TransferERC20 Approve waitForTransaction: %w", err)
	}

	txHash, err := s.erc20Client.Transfer(auth, to, amount)
	if err != nil {
		return fmt.Errorf("TransferERC20 Transfer: %w", err)
	}
	transferLog := &do.TokenTransferLog{
		TokenInfoID:     tokenInfoID,
		FromAddress:     fromAddress.Hex(),
		ToAddress:       toAddress,
		Amount:          amount.Uint64(),
		Status:          do.StatusPending,
		TransactionHash: txHash.Hex(),
		CreateBy:        globalconst.SystemUser,
		CreateAddr:      fromAddress.Hex(),
		CreatedTime:     time.Now(),
	}
	err = do.DefaultTokenTransferLogManager.Create(transferLog)
	if err != nil {
		return fmt.Errorf("create TokenTransferLog: %w", err)
	}

	err = s.waitForTransaction(ctx, txHash)
	if err != nil {
		log.Error("TransferERC20 Transfer waitForTransaction", "err", err)
		err := do.DefaultTokenTransferLogManager.UpdateStatus(transferLog.ID, do.StatusFailed, "", "")
		if err != nil {
			return fmt.Errorf("TransferERC20 Transfer waitForTransaction: %w", err)
		}
		return nil
	}

	err = do.DefaultTokenTransferLogManager.UpdateStatus(transferLog.ID, do.StatusSuccess, "", "")
	if err != nil {
		return fmt.Errorf("TransferERC20 Transfer waitForTransaction: %w", err)
	}
	log.Info("TransferERC20 Transfer success", "txHash", txHash.Hex())
	return nil
}

func (s *BusinessService) waitForTransaction(ctx context.Context, txHash common.Hash) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("TxReceipt timeout")
		default:
			receipt, err := s.ethClient.TxReceiptByTxHash(txHash)
			if err != nil {
				return fmt.Errorf("TxReceipt error: %w", err)
			}

			if receipt != nil {
				if receipt.Status == types.ReceiptStatusSuccessful {
					return nil
				} else {
					return fmt.Errorf("TxReceipt error")
				}
			}
			time.Sleep(5 * time.Second)
		}
	}
}
