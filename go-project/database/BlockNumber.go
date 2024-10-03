package database

import (
	"time"

	"gorm.io/gorm"
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

type BlockInfoRepository struct {
	db *gorm.DB
}

func NewBlockInfoRepository(db *gorm.DB) *BlockInfoRepository {
	return &BlockInfoRepository{db: db}
}

func (r *BlockInfoRepository) Create(block *BlockInfo) error {
	return r.db.Create(block).Error
}

func (r *BlockInfoRepository) Update(block *BlockInfo) error {
	return r.db.Save(block).Error
}

func (r *BlockInfoRepository) Delete(id int64) error {
	return r.db.Delete(&BlockInfo{}, id).Error
}

func (r *BlockInfoRepository) GetByID(id int64) (*BlockInfo, error) {
	var block BlockInfo
	err := r.db.First(&block, id).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *BlockInfoRepository) GetByBlockHash(hash string) (*BlockInfo, error) {
	var block BlockInfo
	err := r.db.Where("block_hash = ?", hash).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *BlockInfoRepository) GetByBlockNumber(number uint64) (*BlockInfo, error) {
	var block BlockInfo
	err := r.db.Where("block_number = ?", number).First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (r *BlockInfoRepository) PageList(page, pageSize int) ([]BlockInfo, int64, error) {
	var blocks []BlockInfo
	var total int64

	err := r.db.Model(&BlockInfo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Order("block_number DESC").Offset(offset).Limit(pageSize).Find(&blocks).Error
	return blocks, total, err
}

func (r *BlockInfoRepository) GetLatestBlock() (*BlockInfo, error) {
	var block BlockInfo
	err := r.db.Order("block_number DESC").First(&block).Error
	if err != nil {
		return nil, err
	}
	return &block, nil
}
