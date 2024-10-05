package do

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const (
	StatusSuccess = "success"
	StatusFailed  = "failed"
	StatusPending = "pending"
)

type TokenTransferLog struct {
	ID              int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TokenInfoID     int       `gorm:"column:token_info_id;not null" json:"token_info_id"`
	WorkflowID      int       `gorm:"column:workflow_id;not null" json:"workflow_id"`
	FromAddress     string    `gorm:"column:from_address;not null;type:VARCHAR(42)" json:"from_address"`
	ToAddress       string    `gorm:"column:to_address;not null;type:VARCHAR(42)" json:"to_address"`
	ContractAddress string    `gorm:"column:contract_address;not null;type:VARCHAR(42)" json:"contract_address"`
	Amount          uint64    `gorm:"column:amount;not null" json:"amount"`
	TransferData    string    `gorm:"column:transfer_data;not null;type:VARCHAR(512)" json:"transfer_data"`
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

func (r *TokenTransferLogManager) Update(log *TokenTransferLog) error {
	return r.db.Model(&TokenTransferLog{}).
		Where("id = ?", log.ID).
		Updates(map[string]interface{}{
			"token_info_id":    log.TokenInfoID,
			"workflow_id":      log.WorkflowID,
			"from_address":     log.FromAddress,
			"to_address":       log.ToAddress,
			"contract_address": log.ContractAddress,
			"amount":           log.Amount,
			"transfer_data":    log.TransferData,
			"status":           log.Status,
			"retry_count":      log.RetryCount,
			"transaction_hash": log.TransactionHash,
			"updated_by":       log.UpdatedBy,
			"updated_addr":     log.UpdatedAddr,
			"updated_time":     time.Now(),
		}).Error
}

func (r *TokenTransferLogManager) Create(log *TokenTransferLog) error {
	return r.db.Create(log).Error
}

func (r *TokenTransferLogManager) GetPendingTokenTransferLogs() ([]TokenTransferLog, error) {
	var logs []TokenTransferLog

	err := r.db.Where("status = ? AND retry_count <= ? and transaction_hash = ''", "pending", 3).
		Limit(10).
		Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("GetPendingTokenTransferLogs err: %w", err)
	}

	return logs, nil
}

func (r *TokenTransferLogManager) GetByTxHashAndAddresses(txHash, from, to string) (*TokenTransferLog, error) {
	var log TokenTransferLog
	err := r.db.Where("transaction_hash = ? AND from_address = ? and status = 'pending'", txHash, from).
		First(&log).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByTxHashAndAddresses err: %w", err)
	}
	return &log, nil
}
