package do

import (
	"gorm.io/gorm"
	"time"
)

type TokenInfo struct {
	ID              int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TokenName       string    `gorm:"column:token_name;not null;type:VARCHAR(100)" json:"token_name"`
	TokenSymbol     string    `gorm:"column:token_symbol;not null;type:VARCHAR(64)" json:"token_symbol"`
	ContractAddress string    `gorm:"column:contract_address;not null;type:VARCHAR(64)" json:"contract_address"`
	Decimals        int       `gorm:"column:decimals;not null;default:18" json:"decimals"`
	CreateBy        string    `gorm:"column:create_by;not null;type:VARCHAR(64)" json:"create_by"`
	CreateAddr      string    `gorm:"column:create_addr;not null;type:VARCHAR(64)" json:"create_addr"`
	CreatedTime     time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
	UpdatedBy       string    `gorm:"column:updated_by;type:VARCHAR(64)" json:"updated_by"`
	UpdatedAddr     string    `gorm:"column:updated_addr;type:VARCHAR(64)" json:"updated_addr"`
	UpdatedTime     time.Time `gorm:"column:updated_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_time"`
}

func (TokenInfo) TableName() string {
	return "token_info"
}

type TokenInfoManager struct {
	db *gorm.DB
}

func NewTokenInfoManager(db *gorm.DB) *TokenInfoManager {
	return &TokenInfoManager{db: db}
}
