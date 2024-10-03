package service

import (
	"fmt"
	"go-project/main/log"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-project/business/workflow/do"
	"go-project/business/workflow/dto"
	"go-project/common/types"
)

type Service struct {
	logger *log.ZapLogger
	db     *gorm.DB
}

func NewService(logger *log.ZapLogger, db *gorm.DB) *Service {
	return &Service{
		logger: logger,
		db:     db,
	}
}

func (service *Service) CreateWorkFlowService(dto *dto.WorkflowInfoCreateDTO) (*do.WorkFlowInfo, error) {
	newWorkflow := &do.WorkFlowInfo{
		WorkflowName: dto.WorkflowName,
		ToAddr:       dto.ToAddr,
		Description:  dto.Description,
		CreateBy:     dto.ToAddr,
		CreateAddr:   dto.ToAddr,
		CreatedTime:  time.Now(),
	}
	workflowManager := do.NewWorkFlowInfoManager(service.db)
	err := workflowManager.Create(newWorkflow)
	if err != nil {
		return nil, fmt.Errorf("CreateWorkFlow error: %w", err)
	}
	return newWorkflow, nil
}

func (service *Service) PageWorkFlowList(req types.GenericPageReq[dto.WorkflowInfoCreateDTO]) (*types.GenericPageResp[do.WorkFlowInfo], error) {
	if req.PageNum == 0 {
		req.PageNum = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	resp := &types.GenericPageResp[do.WorkFlowInfo]{
		PageResp: types.PageResp{
			PageNum:  req.PageNum,
			PageSize: req.PageSize,
		},
	}

	offset := (resp.PageNum - 1) * resp.PageSize

	workflowManager := do.NewWorkFlowInfoManager(service.db)
	list, err := workflowManager.Page(offset, resp.PageSize)
	if err != nil {
		service.logger.Error("PageWorkFlowList Page", zap.Any("err", err))
		return nil, err
	}

	total, err := workflowManager.Count()
	if err != nil {
		service.logger.Error("PageWorkFlowList Count", zap.Any("err", err))
		return nil, err
	}

	resp.List = list
	resp.TotalPage = (total + resp.PageSize - 1) / resp.PageSize
	return resp, nil
}

func (service *Service) ApproveWorkFlow(input *dto.WorkFlowApprovalDTO) error {
	workflowManager := do.NewWorkFlowInfoManager(service.db)
	workflow, err := workflowManager.GetByID(input.WorkflowID)
	if err != nil {
		return fmt.Errorf("getById error: %w", err)
	}
	if workflow == nil {
		return fmt.Errorf("getById is null")
	}

	approve := &do.WorkFlowApprove{
		WorkflowID:  input.WorkflowID,
		ApproveAddr: input.ApproverAddr,
		Status:      input.ApprovalStatus,
		ApproveTime: time.Now(),
		CreateBy:    input.ApproverAddr,
		CreateAddr:  input.ApproverAddr,
		CreatedTime: time.Now(),
	}

	workflowApproveManager := do.NewWorkFlowApproveManager(service.db)
	err = workflowApproveManager.Create(approve)
	if err != nil {
		service.logger.Error("PageWorkFlowList Count", zap.Any("err", err))
		return err
	}

	return nil
}
