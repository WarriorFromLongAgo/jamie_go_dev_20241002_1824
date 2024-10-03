package database

import (
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

type ManagementRepository struct {
	db *gorm.DB
}

func NewManagementRepository(db *gorm.DB) *ManagementRepository {
	return &ManagementRepository{db: db}
}

func (r *ManagementRepository) Create(management *Management) error {
	return r.db.Create(management).Error
}

func (r *ManagementRepository) Update(management *Management) error {
	return r.db.Save(management).Error
}

func (r *ManagementRepository) Delete(id int) error {
	return r.db.Delete(&Management{}, id).Error
}

func (r *ManagementRepository) GetByID(id int) (*Management, error) {
	var management Management
	err := r.db.First(&management, id).Error
	if err != nil {
		return nil, err
	}
	return &management, nil
}

func (r *ManagementRepository) GetByAddr(addr string) (*Management, error) {
	var management Management
	err := r.db.Where("addr = ?", addr).First(&management).Error
	if err != nil {
		return nil, err
	}
	return &management, nil
}

func (r *ManagementRepository) PageList(page, pageSize int) ([]Management, int64, error) {
	var managements []Management
	var total int64

	err := r.db.Model(&Management{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Order("created_time DESC").Offset(offset).Limit(pageSize).Find(&managements).Error
	return managements, total, err
}

func (r *ManagementRepository) GetByPermissionLevel(level string) ([]Management, error) {
	var managements []Management
	err := r.db.Where("permission_level = ?", level).Find(&managements).Error
	return managements, err
}

func (r *ManagementRepository) UpdatePermissionLevel(id int, level string, updatedBy, updatedAddr string) error {
	return r.db.Model(&Management{}).Where("id = ?", id).Updates(map[string]interface{}{
		"permission_level": level,
		"updated_by":       updatedBy,
		"updated_addr":     updatedAddr,
	}).Error
}
