package do

import (
	"gorm.io/gorm"
	"time"
)

type WorkFlowApprove struct {
	ID          int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	WorkflowID  int       `gorm:"column:workflow_id" json:"workflow_id"`
	ApproveAddr string    `gorm:"column:approve_addr;not null;type:VARCHAR(64)" json:"approve_addr"`
	Status      string    `gorm:"column:status;type:ENUM('approved','rejected');default:rejected" json:"status"`
	ApproveTime time.Time `gorm:"column:approve_time" json:"approve_time"`
	CreateBy    string    `gorm:"column:create_by;not null;type:VARCHAR(64)" json:"create_by"`
	CreateAddr  string    `gorm:"column:create_addr;not null;type:VARCHAR(64)" json:"create_addr"`
	CreatedTime time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
	UpdatedBy   string    `gorm:"column:updated_by;type:VARCHAR(64)" json:"updated_by"`
	UpdatedAddr string    `gorm:"column:updated_addr;type:VARCHAR(64)" json:"updated_addr"`
	UpdatedTime time.Time `gorm:"column:updated_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_time"`
}

func (WorkFlowApprove) TableName() string {
	return "workflow_approve"
}

type WorkFlowApproveManager struct {
	db *gorm.DB
}

func NewWorkFlowApproveManager(db *gorm.DB) *WorkFlowApproveManager {
	return &WorkFlowApproveManager{db: db}
}

func (m *WorkFlowApproveManager) Create(approve *WorkFlowApprove) error {
	return m.db.Create(approve).Error
}

func (m *WorkFlowApproveManager) CountUniqueApprovedAddresses(workflowID int) (int64, error) {
	var count int64
	err := m.db.Model(&WorkFlowApprove{}).
		Where("workflow_id = ? AND status = ?", workflowID, "approved").
		Distinct("approve_addr").
		Count(&count).Error
	return count, err
}
