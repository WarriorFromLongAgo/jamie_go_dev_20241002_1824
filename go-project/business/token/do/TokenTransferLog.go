package do

import (
	"gorm.io/gorm"
	"time"
)

const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
	StatusPending = "pending"
)

type TokenTransferLog struct {
	ID              int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TokenInfoID     int       `gorm:"column:token_info_id;not null" json:"token_info_id"`
	FromAddress     string    `gorm:"column:from_address;not null;type:VARCHAR(42)" json:"from_address"`
	ToAddress       string    `gorm:"column:to_address;not null;type:VARCHAR(42)" json:"to_address"`
	Amount          uint64    `gorm:"column:amount;not null" json:"amount"`
	Status          string    `gorm:"column:status;not null;type:ENUM('failed','success','pending');default:pending" json:"status"`
	RetryCount      int       `gorm:"column:retry_count;not null;default:0" json:"retry_count"`
	TransactionHash string    `gorm:"column:transaction_hash;not null;type:VARCHAR(66)" json:"transaction_hash"`
	CreateBy        string    `gorm:"column:create_by;not null;type:VARCHAR(64)" json:"create_by"`
	CreateAddr      string    `gorm:"column:create_addr;not null;type:VARCHAR(64)" json:"create_addr"`
	CreatedTime     time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
	UpdatedBy       string    `gorm:"column:updated_by;type:VARCHAR(64)" json:"updated_by"`
	UpdatedAddr     string    `gorm:"column:updated_addr;type:VARCHAR(64)" json:"updated_addr"`
	UpdatedTime     time.Time `gorm:"column:updated_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_time"`
}

func (TokenTransferLog) TableName() string {
	return "token_transfer_log"
}

type TokenTransferLogManager struct {
	db *gorm.DB
}

func NewTokenTransferLogManager(db *gorm.DB) *TokenTransferLogManager {
	return &TokenTransferLogManager{db: db}
}

func (r *TokenTransferLogManager) UpdateStatus(id int, status string, updatedBy, updatedAddr string) error {
	return r.db.Model(&TokenTransferLog{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"updated_by":   updatedBy,
			"updated_addr": updatedAddr,
			"updated_time": time.Now(),
		}).Error
}

func (r *TokenTransferLogManager) Create(log *TokenTransferLog) error {
	return r.db.Create(log).Error
}
