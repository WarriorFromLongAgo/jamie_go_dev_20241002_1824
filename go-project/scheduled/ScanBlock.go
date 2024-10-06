package scheduled

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
	"gorm.io/gorm"

	do2 "go-project/business/scan/do"
	"go-project/business/token/do"
	"go-project/chain/eth"
	"go-project/main/log"
)

type ScanBlock struct {
	ctx       context.Context
	ethClient eth.EthClient
	db        *gorm.DB
	log       *log.ZapLogger
}

func NewScanBlock(ctx context.Context, client eth.EthClient, db *gorm.DB, log *log.ZapLogger) (*ScanBlock, error) {
	return &ScanBlock{
		ctx:       ctx,
		ethClient: client,
		db:        db,
		log:       log,
	}, nil
}

func (s *ScanBlock) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("ScanBlock done")
			return
		case <-ticker.C:
			err := s.scanBlocks()
			if err != nil {
				fmt.Printf("ScanBlock error: %v\n", err)
			}
		}
	}
}

func (s *ScanBlock) scanBlocks() error {
	blockInfoManager := do2.NewBlockInfoManager(s.db)
	dbLatestBlockNumber, err := blockInfoManager.GetLatestBlockNumber()
	if err != nil {
		return fmt.Errorf("获取最新扫描的区块号失败: %w", err)
	}
	s.log.Info("扫描区块 dbLatestBlockNumber ", zap.Uint64("dbLatestBlockNumber", dbLatestBlockNumber))
	fmt.Printf("扫描区块 dbLatestBlockNumber %d", dbLatestBlockNumber)
	fmt.Println()

	remoteLatestBlock, err := s.ethClient.LatestFinalizedBlockHeader()
	if err != nil {
		return fmt.Errorf("获取最新区块失败: %w", err)
	}
	s.log.Info("远程最新区块 remoteLatestBlock", zap.Uint64("remoteLatestBlock", remoteLatestBlock.Number.Uint64()))
	fmt.Printf("扫描区块 remoteLatestBlock %d", remoteLatestBlock.Number)
	fmt.Println()

	startBlock := new(big.Int).SetUint64(dbLatestBlockNumber)
	if dbLatestBlockNumber > 0 {
		startBlock.Add(startBlock, big.NewInt(1))
	}
	endBlock := new(big.Int).Add(startBlock, big.NewInt(100))
	s.log.Info("扫描区块范围", zap.Uint64("startBlock", startBlock.Uint64()), zap.Uint64("endBlock", endBlock.Uint64()))

	if endBlock.Cmp(remoteLatestBlock.Number) > 0 {
		endBlock = remoteLatestBlock.Number
	}
	if startBlock.Cmp(endBlock) > 0 {
		s.log.Info("没有新区块需要扫描")
		return nil
	}

	headers, err := s.ethClient.BlockHeaderListByRange(startBlock, endBlock)
	if err != nil {
		return fmt.Errorf("获取区块头列表失败: %w", err)
	}

	return s.processBlocksInTransaction(headers)
}

func (s *ScanBlock) processBlocksInTransaction(headers []*types.Header) error {
	if len(headers) == 0 {
		return nil
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		blockInfoManager := do2.NewBlockInfoManager(tx)
		transactionManager := do2.NewTransactionInfoManager(tx)
		tokenTransferLogManager := do.NewTokenTransferLogManager(tx)

		for _, header := range headers {
			if err := s.processBlockHeader(header, blockInfoManager); err != nil {
				return err
			}

			if err := s.processBlockTransactions(header, transactionManager, tokenTransferLogManager); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *ScanBlock) processBlockHeader(header *types.Header, blockInfoManager *do2.BlockInfoManager) error {
	s.log.Info("处理区块头", zap.Uint64("blockNumber", header.Number.Uint64()), zap.String("blockHash", header.Hash().Hex()))

	err := blockInfoManager.Create(&do2.BlockInfo{
		BlockNumber:     header.Number.Uint64(),
		BlockHash:       header.Hash().Hex(),
		BlockParentHash: header.ParentHash.Hex(),
		Timestamp:       time.Unix(int64(header.Time), 0),
	})
	if err != nil {
		return fmt.Errorf("保存区块信息失败 (区块 %d): %w", header.Number, err)
	}

	return nil
}

func (s *ScanBlock) processBlockTransactions(header *types.Header, transactionManager *do2.TransactionInfoManager, tokenTransferLogManager *do.TokenTransferLogManager) error {
	block, err := s.ethClient.BlockByNumberV3(s.ctx, header.Number)
	if err != nil {
		return fmt.Errorf("获取区块失败: %w", err)
	}

	s.log.Info("处理区块交易", zap.Uint64("blockNumber", block.NumberU64()), zap.Int("txCount", len(block.Transactions())))

	if len(block.Transactions()) == 0 {
		return nil
	}

	for _, tx := range block.Transactions() {
		if err := s.processSingleTransaction(block, tx, transactionManager, tokenTransferLogManager); err != nil {
			return err
		}
	}

	return nil
}

func (s *ScanBlock) processSingleTransaction(block *types.Block, tx *types.Transaction, transactionManager *do2.TransactionInfoManager, tokenTransferLogManager *do.TokenTransferLogManager) error {
	txHash := tx.Hash().Hex()
	// 尝试使用不同的方法获取发送者
	var from common.Address
	var err error
	signer := types.LatestSignerForChainID(tx.ChainId())
	from, err = types.Sender(signer, tx)
	if err != nil {
		// 如果仍然失败，记录错误并继续处理其他字段
		s.log.Error("无法获取交易发送者", zap.Error(err), zap.String("txHash", txHash))
		from = common.Address{}
	}

	to := tx.To()
	toAddress := ""
	if to != nil {
		toAddress = to.Hex()
	}

	data := tx.Data()
	if len(data) >= 4 && toAddress == "" {
		if len(data) >= 36 {
			toAddress = common.BytesToAddress(data[4:36]).Hex()
		}
	}

	receipt, err := s.ethClient.TxReceiptByTxHash(tx.Hash())
	if err != nil {
		return fmt.Errorf("获取交易收据失败: %w", err)
	}

	txInfo := &do2.TransactionInfo{
		BlockNumber:      block.NumberU64(),
		BlockHash:        block.Hash().Hex(),
		TxHash:           txHash,
		FromAddress:      from.Hex(),
		ToAddress:        toAddress,
		TokenAddress:     toAddress,
		Value:            tx.Value().String(),
		GasPrice:         tx.GasPrice().String(),
		GasLimit:         tx.Gas(),
		GasUsed:          receipt.GasUsed,
		Nonce:            tx.Nonce(),
		TransactionIndex: uint64(receipt.TransactionIndex),
		Status:           receipt.Status,
		TxType:           tx.Type(),
		Data:             common.Bytes2Hex(tx.Data()),
		CreatedTime:      time.Unix(int64(block.Time()), 0),
	}

	if err := transactionManager.Create(txInfo); err != nil {
		return fmt.Errorf("保存交易信息失败: %w", err)
	}

	if err := s.updateTokenTransferLog(txHash, from.Hex(), toAddress, tokenTransferLogManager); err != nil {
		return fmt.Errorf("更新TokenTransferLog失败: %w", err)
	}

	return nil
}

func (s *ScanBlock) updateTokenTransferLog(txHash, fromAddress, toAddress string, tokenTransferLogManager *do.TokenTransferLogManager) error {
	pendingLog, err := tokenTransferLogManager.GetByTxHashAndAddresses(txHash, fromAddress, toAddress)
	if err != nil {
		return fmt.Errorf("查询TokenTransferLog失败: %w", err)
	}

	if pendingLog != nil && pendingLog.Status == do.StatusPending {
		pendingLog.Status = do.StatusSuccess
		pendingLog.UpdatedTime = time.Now()
		pendingLog.UpdatedBy = "ScanBlock"
		pendingLog.UpdatedAddr = "system"

		if err := tokenTransferLogManager.Update(pendingLog); err != nil {
			return fmt.Errorf("更新TokenTransferLog状态失败: %w", err)
		}

		s.log.Info("TokenTransferLog状态更新为成功", zap.String("txHash", txHash))
	}

	return nil
}
