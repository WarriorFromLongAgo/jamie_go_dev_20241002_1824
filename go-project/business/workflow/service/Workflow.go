package service

import (
	"fmt"
	"go-project/main/app"
	"gorm.io/gorm"
	"time"

	"go.uber.org/zap"

	"go-project/business/workflow/do"
	"go-project/business/workflow/dto"
	"go-project/common/types"
)

type Service struct {
	logger app.Logger
	db     *gorm.DB
}

func NewService(logger app.Logger, db *gorm.DB) *Service {
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
	err := do.DefaultWorkFlowInfoManager.Create(newWorkflow)
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
	list, err := do.DefaultWorkFlowInfoManager.Page(offset, resp.PageSize)
	if err != nil {
		app.Log.Error("PageWorkFlowList Page", zap.Any("err", err))
		return nil, err
	}

	total, err := do.DefaultWorkFlowInfoManager.Count()
	if err != nil {
		app.Log.Error("PageWorkFlowList Count", zap.Any("err", err))
		return nil, err
	}

	resp.List = list
	resp.TotalPage = (total + resp.PageSize - 1) / resp.PageSize
	return resp, nil
}

func (service *Service) ApproveWorkFlow(input *dto.WorkFlowApprovalDTO) error {
	workflow, err := do.DefaultWorkFlowInfoManager.GetByID(input.WorkflowID)
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

	err = do.DefaultWorkFlowApproveManager.Create(approve)
	if err != nil {
		app.Log.Error("PageWorkFlowList Count", zap.Any("err", err))
		return err
	}

	return nil
}
