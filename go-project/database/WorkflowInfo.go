package database

import (
	"time"

	"gorm.io/gorm"
)

type WorkflowInfo struct {
	ID           int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	WorkflowName string    `gorm:"column:workflow_name;not null;type:VARCHAR(128)" json:"workflow_name"`
	ToAddr       string    `gorm:"column:to_addr;not null;type:VARCHAR(64)" json:"to_addr"`
	Description  string    `gorm:"column:description;not null;type:VARCHAR(1024)" json:"description"`
	CreateBy     string    `gorm:"column:create_by;not null;type:VARCHAR(64)" json:"create_by"`
	CreateAddr   string    `gorm:"column:create_addr;not null;type:VARCHAR(64)" json:"create_addr"`
	CreatedTime  time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
	UpdatedBy    string    `gorm:"column:updated_by;type:VARCHAR(64)" json:"updated_by"`
	UpdatedAddr  string    `gorm:"column:updated_addr;type:VARCHAR(64)" json:"updated_addr"`
	UpdatedTime  time.Time `gorm:"column:updated_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_time"`
}

func (WorkflowInfo) TableName() string {
	return "workflow_info"
}

type WorkflowInfoRepository struct {
	db *gorm.DB
}

func NewWorkflowInfoRepository(db *gorm.DB) *WorkflowInfoRepository {
	return &WorkflowInfoRepository{db: db}
}

func (r *WorkflowInfoRepository) Create(workflow *WorkflowInfo) error {
	return r.db.Create(workflow).Error
}

func (r *WorkflowInfoRepository) Update(workflow *WorkflowInfo) error {
	return r.db.Save(workflow).Error
}

func (r *WorkflowInfoRepository) Delete(id int) error {
	return r.db.Delete(&WorkflowInfo{}, id).Error
}

func (r *WorkflowInfoRepository) GetByID(id int) (*WorkflowInfo, error) {
	var workflow WorkflowInfo
	err := r.db.First(&workflow, id).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (r *WorkflowInfoRepository) PageList(page, pageSize int) ([]WorkflowInfo, int64, error) {
	var workflows []WorkflowInfo
	var total int64

	err := r.db.Model(&WorkflowInfo{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Order("created_time DESC").Offset(offset).Limit(pageSize).Find(&workflows).Error
	return workflows, total, err
}
