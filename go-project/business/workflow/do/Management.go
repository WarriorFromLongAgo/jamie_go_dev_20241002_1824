package do

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Management struct {
	ID              int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name            string    `gorm:"column:name;not null;type:VARCHAR(100)" json:"name"`
	PermissionLevel string    `gorm:"column:permission_level;type:ENUM('none','partial','full');default:none" json:"permission_level"`
	Addr            string    `gorm:"column:addr;not null;type:VARCHAR(64)" json:"addr"`
	AnvilInfo       string    `gorm:"column:anvil_info;not null;type:VARCHAR(64)" json:"anvil_info"`
	CreateBy        string    `gorm:"column:create_by;not null;type:VARCHAR(64)" json:"create_by"`
	CreateAddr      string    `gorm:"column:create_addr;not null;type:VARCHAR(64)" json:"create_addr"`
	CreatedTime     time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
	UpdatedBy       string    `gorm:"column:updated_by;type:VARCHAR(64)" json:"updated_by"`
	UpdatedAddr     string    `gorm:"column:updated_addr;type:VARCHAR(64)" json:"updated_addr"`
	UpdatedTime     time.Time `gorm:"column:updated_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_time"`
}

func (Management) TableName() string {
	return "management"
}

type ManagementManager struct {
	db *gorm.DB
}

func NewManagementManager(db *gorm.DB) *ManagementManager {
	return &ManagementManager{db: db}
}

func (m *ManagementManager) HasFullPermission(addr string) (bool, error) {
	var management Management
	result := m.db.Where("addr = ? AND permission_level = ?", addr, "full").First(&management)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
