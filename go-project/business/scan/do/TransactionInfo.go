package do

import (
	"time"
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

type TransactionInfoManager struct{}

var DefaultTransactionInfoManager = &TransactionInfoManager{}
