package scheduled

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-project/business/token/do"
	do2 "go-project/business/workflow/do"
	"go-project/chain/eth"
	globalconst "go-project/common"
	"go-project/main/log"
)

type ProcessingFLow struct {
	ctx         context.Context
	ethClient   eth.EthClient
	erc20Client eth.TestErc20Client
	db          *gorm.DB
	log         *log.ZapLogger
}

func NewProcessingFLow(ctx context.Context, client eth.EthClient, erc20Client eth.TestErc20Client, db *gorm.DB, log *log.ZapLogger) (*ProcessingFLow, error) {
	return &ProcessingFLow{
		ctx:         ctx,
		ethClient:   client,
		erc20Client: erc20Client,
		db:          db,
		log:         log,
	}, nil
}

func (s *ProcessingFLow) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("ProcessingFLow done")
			return
		case <-ticker.C:
			fmt.Println("ProcessingFLow start")
			err := s.processingFLow()
			if err != nil {
				fmt.Printf("ProcessingFLow error: %v\n", err)
			}
		}
	}
}

func (s *ProcessingFLow) processingFLow() error {
	tokenTransferLogManager := do.NewTokenTransferLogManager(s.db)
	pendingLogList, err := tokenTransferLogManager.GetPendingTokenTransferLogs()
	if err != nil {
		s.log.Error("processingFLow GetPendingTokenTransferLogs", zap.Error(err))
		return err
	}
	if len(pendingLogList) <= 0 {
		s.log.Info("processingFLow pendingLogList is nil")
		return nil
	}
	pendingLogListJson, _ := json.Marshal(pendingLogList)
	s.log.Info("current deal pendingLogList", zap.Any("pendingLogList", string(pendingLogListJson)))

	tokenInfoManager := do.NewTokenInfoManager(s.db)
	tokenInfo, err := tokenInfoManager.GetByID(1)
	if err != nil {
		s.log.Error("processingFLow GetByID", zap.Error(err))
		return err
	}

	businessService := eth.NewEthBusinessService(s.ethClient, s.erc20Client, s.log)

	for _, pendingLog := range pendingLogList {
		workflowManager := do2.NewWorkFlowInfoManager(s.db)
		workflow, err := workflowManager.GetByID(pendingLog.WorkflowID)
		if err != nil {
			s.log.Error("获取工作流信息失败", zap.Error(err), zap.Int("WorkflowID", pendingLog.WorkflowID))
			continue
		}

		privateKey, err := crypto.HexToECDSA(globalconst.OWNER_PRV_KEY)
		if err != nil {
			s.log.Error("解析私钥失败", zap.Error(err))
			continue
		}
		fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

		printERC20Balance(s.ctx, s.erc20Client, fromAddress, "From (after)")
		printERC20Balance(s.ctx, s.erc20Client, common.HexToAddress(workflow.ToAddr), "To (after)")

		txHash, transferData, err := businessService.TransferERC20(
			s.ctx,
			privateKey,
			fromAddress.Hex(),
			workflow.ToAddr,
			tokenInfo.ContractAddress,
			new(big.Int).SetUint64(123456),
		)

		if err != nil {
			s.log.Error("ERC20转账失败", zap.Error(err), zap.Int("LogID", pendingLog.ID))
			pendingLog.RetryCount++
			pendingLog.Status = do.StatusPending
		} else {
			s.log.Info("ERC20转账成功", zap.Int("LogID", pendingLog.ID), zap.String("TxHash", txHash))
			pendingLog.Status = do.StatusPending
			pendingLog.TransactionHash = txHash
		}

		pendingLog.TransferData = hexutil.Encode(transferData)
		pendingLog.FromAddress = fromAddress.Hex()
		pendingLog.ToAddress = workflow.ToAddr
		pendingLog.ContractAddress = tokenInfo.ContractAddress
		pendingLog.UpdatedBy = fromAddress.Hex()
		pendingLog.UpdatedAddr = fromAddress.Hex()
		pendingLog.UpdatedTime = time.Now()

		err = tokenTransferLogManager.Update(&pendingLog)
		if err != nil {
			s.log.Error("更新转账日志状态失败", zap.Error(err), zap.Int("LogID", pendingLog.ID))
		}

		printERC20Balance(s.ctx, s.erc20Client, fromAddress, "From (after)")
		printERC20Balance(s.ctx, s.erc20Client, common.HexToAddress(workflow.ToAddr), "To (after)")
	}

	return nil
}

func printERC20Balance(ctx context.Context, client eth.TestErc20Client, address common.Address, label string) {
	balance, err := client.BalanceOf(ctx, address)
	if err != nil {
		fmt.Printf("Failed to get balance for %s address: %v", label, err)
		fmt.Println()
	}
	fmt.Printf("%s address balance: %s", label, balance.String())
	fmt.Println()
}
