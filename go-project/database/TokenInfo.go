package database

import (
	"time"

	"gorm.io/gorm"
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

type TokenInfoRepository struct {
	db *gorm.DB
}

func NewTokenInfoRepository(db *gorm.DB) *TokenInfoRepository {
	return &TokenInfoRepository{db: db}
}

func (r *TokenInfoRepository) Create(token *TokenInfo) error {
	return r.db.Create(token).Error
}

func (r *TokenInfoRepository) Update(token *TokenInfo) error {
	return r.db.Save(token).Error
}

func (r *TokenInfoRepository) Delete(id int) error {
	return r.db.Delete(&TokenInfo{}, id).Error
}

func (r *TokenInfoRepository) GetByID(id int) (*TokenInfo, error) {
	var token TokenInfo
	err := r.db.First(&token, id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenInfoRepository) GetByContractAddress(address string) (*TokenInfo, error) {
	var token TokenInfo
	err := r.db.Where("contract_address = ?", address).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenInfoRepository) GetBySymbol(symbol string) (*TokenInfo, error) {
	var token TokenInfo
	err := r.db.Where("token_symbol = ?", symbol).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenInfoRepository) List(page, pageSize int) ([]TokenInfo, int64, error) {
	var tokens []TokenInfo
	var total int64

	err := r.db.Model(&TokenInfo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Order("created_time DESC").Offset(offset).Limit(pageSize).Find(&tokens).Error
	return tokens, total, err
}

func (r *TokenInfoRepository) Search(keyword string) ([]TokenInfo, error) {
	var tokens []TokenInfo
	err := r.db.Where("token_name LIKE ? OR token_symbol LIKE ? OR contract_address LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%").
		Find(&tokens).Error
	return tokens, err
}
