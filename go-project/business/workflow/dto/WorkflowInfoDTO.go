package dto

type WorkflowInfoCreateDTO struct {
	WorkflowName string `json:"workflow_name" binding:"required,max=128"`
	ToAddr       string `json:"to_addr" binding:"required,max=64"`
	Description  string `json:"description" binding:"max=1024"`
}

type WorkFlowApprovalDTO struct {
	WorkflowID     int    `json:"workflow_id" binding:"required"`
	ApprovalStatus string `json:"approval_status" binding:"required,oneof=approved rejected"`
	ApproverID     string `json:"approver_id" binding:"required"`
	ApproverAddr   string `json:"approver_addr" binding:"required"`
}
