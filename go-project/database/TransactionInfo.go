package database

import (
	"time"

	"gorm.io/gorm"
)

type TransactionInfo struct {
	ID               int64     `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	BlockHash        string    `gorm:"column:block_hash;not null;type:VARCHAR(128)" json:"block_hash"`
	BlockNumber      uint64    `gorm:"column:block_number;not null;index" json:"block_number"`
	TxHash           string    `gorm:"column:tx_hash;not null;uniqueIndex;type:VARCHAR(128)" json:"tx_hash"`
	FromAddress      string    `gorm:"column:from_address;not null;index;type:VARCHAR(64)" json:"from_address"`
	ToAddress        string    `gorm:"column:to_address;not null;index;type:VARCHAR(128)" json:"to_address"`
	TokenAddress     string    `gorm:"column:token_address;not null;type:VARCHAR(128)" json:"token_address"`
	GasFee           int64     `gorm:"column:gas_fee;not null" json:"gas_fee"`
	Amount           int64     `gorm:"column:amount;not null" json:"amount"`
	Status           int16     `gorm:"column:status;not null;default:0" json:"status"`
	TransactionIndex int64     `gorm:"column:transaction_index;not null" json:"transaction_index"`
	TxType           int16     `gorm:"column:tx_type;not null;default:0" json:"tx_type"`
	CreatedTime      time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
}

func (TransactionInfo) TableName() string {
	return "transaction_info"
}

type TransactionInfoRepository struct {
	db *gorm.DB
}

func NewTransactionInfoRepository(db *gorm.DB) *TransactionInfoRepository {
	return &TransactionInfoRepository{db: db}
}

func (r *TransactionInfoRepository) Create(tx *TransactionInfo) error {
	return r.db.Create(tx).Error
}

func (r *TransactionInfoRepository) Update(tx *TransactionInfo) error {
	return r.db.Save(tx).Error
}

func (r *TransactionInfoRepository) Delete(id int64) error {
	return r.db.Delete(&TransactionInfo{}, id).Error
}

func (r *TransactionInfoRepository) GetByID(id int64) (*TransactionInfo, error) {
	var tx TransactionInfo
	err := r.db.First(&tx, id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionInfoRepository) GetByTxHash(hash string) (*TransactionInfo, error) {
	var tx TransactionInfo
	err := r.db.Where("tx_hash = ?", hash).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *TransactionInfoRepository) GetByBlockHash(blockHash string) ([]TransactionInfo, error) {
	var txs []TransactionInfo
	err := r.db.Where("block_hash = ?", blockHash).Find(&txs).Error
	return txs, err
}

func (r *TransactionInfoRepository) PageList(page, pageSize int) ([]TransactionInfo, int64, error) {
	var txs []TransactionInfo
	var total int64

	err := r.db.Model(&TransactionInfo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Offset(offset).Limit(pageSize).Find(&txs).Error
	return txs, total, err
}
