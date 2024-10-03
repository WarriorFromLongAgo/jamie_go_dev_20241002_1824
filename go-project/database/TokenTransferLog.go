package database

import (
	"time"

	"gorm.io/gorm"
)

type TokenTransferLog struct {
	ID              int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TokenInfoID     int       `gorm:"column:token_info_id;not null" json:"token_info_id"`
	FromAddress     string    `gorm:"column:from_address;not null;type:VARCHAR(42)" json:"from_address"`
	ToAddress       string    `gorm:"column:to_address;not null;type:VARCHAR(42)" json:"to_address"`
	Amount          uint64    `gorm:"column:amount;not null" json:"amount"`
	Status          string    `gorm:"column:status;not null;type:ENUM('pending','success','failed');default:pending" json:"status"`
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

type TokenTransferLogRepository struct {
	db *gorm.DB
}

func NewTokenTransferLogRepository(db *gorm.DB) *TokenTransferLogRepository {
	return &TokenTransferLogRepository{db: db}
}

func (r *TokenTransferLogRepository) Create(log *TokenTransferLog) error {
	return r.db.Create(log).Error
}

func (r *TokenTransferLogRepository) Update(log *TokenTransferLog) error {
	return r.db.Save(log).Error
}

func (r *TokenTransferLogRepository) Delete(id int) error {
	return r.db.Delete(&TokenTransferLog{}, id).Error
}

func (r *TokenTransferLogRepository) GetByID(id int) (*TokenTransferLog, error) {
	var log TokenTransferLog
	err := r.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *TokenTransferLogRepository) GetByTransactionHash(hash string) (*TokenTransferLog, error) {
	var log TokenTransferLog
	err := r.db.Where("transaction_hash = ?", hash).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *TokenTransferLogRepository) PageList(page, pageSize int) ([]TokenTransferLog, int64, error) {
	var logs []TokenTransferLog
	var total int64

	err := r.db.Model(&TokenTransferLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Order("created_time DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

func (r *TokenTransferLogRepository) GetPendingTransfers(limit int) ([]TokenTransferLog, error) {
	var logs []TokenTransferLog
	err := r.db.Where("status = ?", "pending").Order("created_time ASC").Limit(limit).Find(&logs).Error
	return logs, err
}

func (r *TokenTransferLogRepository) UpdateStatus(id int, status string, updatedBy, updatedAddr string) error {
	return r.db.Model(&TokenTransferLog{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"updated_by":   updatedBy,
		"updated_addr": updatedAddr,
	}).Error
}
