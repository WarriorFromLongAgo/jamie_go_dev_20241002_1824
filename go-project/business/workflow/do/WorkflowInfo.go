package do

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	WorkFlowStatusPending  = "pending"
	WorkFlowStatusApproved = "approved"
	WorkFlowStatusRejected = "rejected"
)

type WorkFlowInfo struct {
	ID           int       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	WorkflowName string    `gorm:"column:workflow_name;not null;type:VARCHAR(128)" json:"workflow_name"`
	ToAddr       string    `gorm:"column:to_addr;not null;type:VARCHAR(64)" json:"to_addr"`
	TokenInfoID  int       `gorm:"column:token_info_id;not null" json:"token_info_id"`
	Description  string    `gorm:"column:description;not null;type:VARCHAR(1024)" json:"description"`
	Status       string    `gorm:"column:status;type:ENUM('pending','approved','rejected');default:pending" json:"status"`
	CreateBy     string    `gorm:"column:create_by;not null;type:VARCHAR(64)" json:"create_by"`
	CreateAddr   string    `gorm:"column:create_addr;not null;type:VARCHAR(64)" json:"create_addr"`
	CreatedTime  time.Time `gorm:"column:created_time;default:CURRENT_TIMESTAMP" json:"created_time"`
	UpdatedBy    string    `gorm:"column:updated_by;type:VARCHAR(64)" json:"updated_by"`
	UpdatedAddr  string    `gorm:"column:updated_addr;type:VARCHAR(64)" json:"updated_addr"`
	UpdatedTime  time.Time `gorm:"column:updated_time;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_time"`
}

func (WorkFlowInfo) TableName() string {
	return "workflow_info"
}

type WorkFlowInfoManager struct {
	db *gorm.DB
}

func NewWorkFlowInfoManager(db *gorm.DB) *WorkFlowInfoManager {
	return &WorkFlowInfoManager{db: db}
}

func (m *WorkFlowInfoManager) Create(info *WorkFlowInfo) error {
	return m.db.Create(info).Error
}

func (m *WorkFlowInfoManager) GetByID(id int) (*WorkFlowInfo, error) {
	var info WorkFlowInfo
	err := m.db.Where("id = ?", id).First(&info).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &info, nil
}

func (m *WorkFlowInfoManager) Page(offset, limit uint64) ([]WorkFlowInfo, error) {
	var infos []WorkFlowInfo
	err := m.db.Offset(int(offset)).Limit(int(limit)).Find(&infos).Error
	return infos, err
}

func (m *WorkFlowInfoManager) Count() (uint64, error) {
	var count int64
	err := m.db.Model(&WorkFlowInfo{}).Count(&count).Error
	return uint64(count), err
}

func (m *WorkFlowInfoManager) Update(workflow *WorkFlowInfo) error {
	return m.db.Save(workflow).Error
}

//func (m *WorkFlowInfoManager) GetPendingWorkflows(limit int) ([]*WorkFlowInfo, error) {
//	var workflows []*WorkFlowInfo
//	err := m.db.Where("status = ?", WorkFlowStatusApproved).
//		Joins("LEFT JOIN token_transfer_log ON workflow_info.id = token_transfer_log.workflow_id").
//		Where("token_transfer_log.id IS NULL").
//		Limit(limit).
//		Find(&workflows).Error
//	return workflows, err
//}
