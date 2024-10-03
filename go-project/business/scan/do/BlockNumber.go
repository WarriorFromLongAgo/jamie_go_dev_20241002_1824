package do

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"go-project/main/app"
)

type BlockInfo struct {
	ID              int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	BlockHash       string    `gorm:"column:block_hash;not null;uniqueIndex;type:VARCHAR(128)" json:"block_hash"`
	BlockParentHash string    `gorm:"column:block_parent_hash;not null;type:VARCHAR(128)" json:"block_parent_hash"`
	BlockNumber     uint64    `gorm:"column:block_number;not null;uniqueIndex" json:"block_number"`
	Timestamp       time.Time `gorm:"column:timestamp;not null" json:"timestamp"`
	RlpBytes        string    `gorm:"column:rlp_bytes;not null;type:VARCHAR(128)" json:"rlp_bytes"`
	CreatedTime     time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
}

func (BlockInfo) TableName() string {
	return "block_info"
}

type BlockInfoManager struct{}

var DefaultBlockInfoManager = &BlockInfoManager{}

func (bm *BlockInfoManager) Create(block *BlockInfo) error {
	return app.DB.Create(block).Error
}

func (bm *BlockInfoManager) GetLatestBlock() (*BlockInfo, error) {
	var block BlockInfo
	var result = app.DB.Order("block_number DESC").First(&block)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &block, nil
}

func (bm *BlockInfoManager) GetLatestBlockNumber() (uint64, error) {
	var block BlockInfo
	var result = app.DB.Order("block_number DESC").First(&block)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		return 0, result.Error
	}
	return block.BlockNumber, nil
}
