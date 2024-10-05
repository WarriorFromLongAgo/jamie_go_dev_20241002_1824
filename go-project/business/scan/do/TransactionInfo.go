package do

import (
	"gorm.io/gorm"
	"time"
)

type TransactionInfo struct {
	ID               uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	BlockHash        string    `gorm:"column:block_hash;type:varchar(128);not null"`
	BlockNumber      uint64    `gorm:"column:block_number;type:bigint unsigned;not null"`
	TxHash           string    `gorm:"column:tx_hash;type:varchar(128);not null;uniqueIndex"`
	FromAddress      string    `gorm:"column:from_address;type:varchar(64);not null;index"`
	ToAddress        string    `gorm:"column:to_address;type:varchar(128);index"`
	TokenAddress     string    `gorm:"column:token_address;type:varchar(128);index"`
	Value            string    `gorm:"column:value;type:varchar(128);not null"`
	GasPrice         string    `gorm:"column:gas_price;type:varchar(128);not null"`
	GasLimit         uint64    `gorm:"column:gas_limit;type:bigint unsigned;not null"`
	GasUsed          uint64    `gorm:"column:gas_used;type:bigint unsigned"`
	Nonce            uint64    `gorm:"column:nonce;type:bigint unsigned;not null"`
	TransactionIndex uint64    `gorm:"column:transaction_index;type:bigint unsigned;not null"`
	Status           uint64    `gorm:"column:status;type:bigint unsigned"`
	TxType           uint8     `gorm:"column:tx_type;type:tinyint unsigned;not null"`
	Data             string    `gorm:"column:data;type:text"`
	CreatedTime      time.Time `gorm:"column:created_time;type:timestamp;default:CURRENT_TIMESTAMP"`
}

func (TransactionInfo) TableName() string {
	return "transaction_info"
}

type TransactionInfoManager struct {
	db *gorm.DB
}

func NewTransactionInfoManager(db *gorm.DB) *TransactionInfoManager {
	return &TransactionInfoManager{db: db}
}

func (m *TransactionInfoManager) Create(info *TransactionInfo) error {
	return m.db.Create(info).Error
}
