package synchronizer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	do2 "go-project/business/scan/do"
	"go-project/chain/eth"
)

type Synchronizer struct {
	ctx       context.Context
	ethClient eth.EthClient
}

func NewSynchronizer(client eth.EthClient, ctx context.Context) (*Synchronizer, error) {
	return &Synchronizer{
		ctx:       ctx,
		ethClient: client,
	}, nil
}

func (s *Synchronizer) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			fmt.Println("区块扫描任务停止")
			return
		case <-ticker.C:
			err := s.scanBlocks()
			if err != nil {
				fmt.Printf("扫描区块出错: %v\n", err)
			}
		}
	}
}

func (s *Synchronizer) scanBlocks() error {
	dbLatestBlockNumber, err := do2.DefaultBlockInfoManager.GetLatestBlockNumber()
	if err != nil {
		return fmt.Errorf("获取最新扫描的区块号失败: %w", err)
	}
	remoteLatestBlock, err := s.ethClient.LatestFinalizedBlockHeader()
	if err != nil {
		return fmt.Errorf("获取最新区块失败: %w", err)
	}

	startBlock := new(big.Int).SetUint64(dbLatestBlockNumber)
	endBlock := new(big.Int).Add(startBlock, big.NewInt(100))

	if endBlock.Cmp(remoteLatestBlock.Number) > 0 {
		endBlock = remoteLatestBlock.Number
	}
	if startBlock.Cmp(endBlock) > 0 {
		return nil
	}
	headers, err := s.ethClient.BlockHeaderListByRange(startBlock, endBlock)
	if err != nil {
		return fmt.Errorf("获取区块头列表失败: %w", err)
	}
	err = s.processBlockHeaders(headers)
	if err != nil {
		return fmt.Errorf("处理区块头失败: %w", err)
	}

	fmt.Printf("成功扫描区块 %d 到 %d\n", startBlock, endBlock)
	return nil
}

func (s *Synchronizer) processBlockHeaders(headers []types.Header) error {
	for _, header := range headers {

		err := do2.DefaultBlockInfoManager.Create(&do2.BlockInfo{
			BlockNumber: header.Number.Uint64(),
			BlockHash:   header.Hash().Hex(),
			Timestamp:   time.Unix(int64(header.Time), 0),
		})
		if err != nil {
			return fmt.Errorf("保存区块信息失败 (区块 %d): %w", header.Number, err)
		}
		fmt.Printf("处理区块 %d, 哈希: %s\n", header.Number, header.Hash().Hex())
	}
	return nil
}
