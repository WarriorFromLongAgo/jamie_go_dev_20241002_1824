package database

import (
	"time"

	"gorm.io/gorm"
)

type WorkflowApprove struct {
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

func (WorkflowApprove) TableName() string {
	return "workflow_approve"
}

type WorkflowApproveRepository struct {
	db *gorm.DB
}

func NewWorkflowApproveRepository(db *gorm.DB) *WorkflowApproveRepository {
	return &WorkflowApproveRepository{db: db}
}

func (r *WorkflowApproveRepository) Create(approval *WorkflowApprove) error {
	return r.db.Create(approval).Error
}

func (r *WorkflowApproveRepository) Update(approval *WorkflowApprove) error {
	return r.db.Save(approval).Error
}

func (r *WorkflowApproveRepository) GetByID(id int) (*WorkflowApprove, error) {
	var approval WorkflowApprove
	err := r.db.First(&approval, id).Error
	if err != nil {
		return nil, err
	}
	return &approval, nil
}

func (r *WorkflowApproveRepository) GetByWorkflowID(workflowID int) ([]WorkflowApprove, error) {
	var approvals []WorkflowApprove
	err := r.db.Where("workflow_id = ?", workflowID).Find(&approvals).Error
	return approvals, err
}

func (r *WorkflowApproveRepository) UpdateStatus(id int, status string, approveTime time.Time, updatedBy, updatedAddr string) error {
	return r.db.Model(&WorkflowApprove{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"approve_time": approveTime,
		"updated_by":   updatedBy,
		"updated_addr": updatedAddr,
	}).Error
}

//func (r *WorkflowApproveRepository) GetPendingApprovals() ([]WorkflowApprove, error) {
//	var approvals []WorkflowApprove
//	err := r.db.Where("status = ?", "rejected").Find(&approvals).Error
//	return approvals, err
//}
//
//func (r *WorkflowApproveRepository) GetApprovalsByApproveAddr(approveAddr string) ([]WorkflowApprove, error) {
//	var approvals []WorkflowApprove
//	err := r.db.Where("approve_addr = ?", approveAddr).Find(&approvals).Error
//	return approvals, err
//}
//
//func (r *WorkflowApproveRepository) GetLatestApprovalByWorkflowID(workflowID int) (*WorkflowApprove, error) {
//	var approval WorkflowApprove
//	err := r.db.Where("workflow_id = ?", workflowID).Order("created_time DESC").First(&approval).Error
//	if err != nil {
//		return nil, err
//	}
//	return &approval, nil
//}
