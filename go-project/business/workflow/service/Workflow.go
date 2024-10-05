package service

import (
	"fmt"
	do2 "go-project/business/token/do"
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
	var newWorkflow *do.WorkFlowInfo

	err := service.db.Transaction(func(tx *gorm.DB) error {
		managementManager := do.NewManagementManager(tx)
		hasFullPermission, err := managementManager.HasFullPermission(dto.ToAddr)
		if err != nil {
			return fmt.Errorf("check permission error: %w", err)
		}

		status := do.WorkFlowStatusPending
		if hasFullPermission {
			status = do.WorkFlowStatusApproved
		}

		newWorkflow = &do.WorkFlowInfo{
			WorkflowName: dto.WorkflowName,
			ToAddr:       dto.ToAddr,
			TokenInfoID:  1,
			Description:  dto.Description,
			Status:       status,
			CreateBy:     dto.ToAddr,
			CreateAddr:   dto.ToAddr,
			CreatedTime:  time.Now(),
		}

		workflowManager := do.NewWorkFlowInfoManager(tx)
		err = workflowManager.Create(newWorkflow)
		if err != nil {
			return fmt.Errorf("CreateWorkFlow create error: %w", err)
		}

		if status == do.WorkFlowStatusApproved {
			tokenInfoManager := do2.NewTokenInfoManager(tx)
			tokenInfo, err := tokenInfoManager.GetByID(1)
			if err != nil {
				service.logger.Error("CreateWorkFlow tokenInfoManager GetByID", zap.Error(err))
				return err
			}

			tokenTransferLog := &do2.TokenTransferLog{
				TokenInfoID:     newWorkflow.TokenInfoID,
				WorkflowID:      newWorkflow.ID,
				FromAddress:     "0x0",
				ToAddress:       newWorkflow.ToAddr,
				ContractAddress: tokenInfo.ContractAddress,
				Amount:          0,
				TransferData:    "",
				Status:          do2.StatusPending,
				RetryCount:      0,
				TransactionHash: "",
				CreateBy:        newWorkflow.CreateBy,
				CreateAddr:      newWorkflow.CreateAddr,
				CreatedTime:     time.Now(),
			}

			tokenTransferLogManager := do2.NewTokenTransferLogManager(tx)
			err = tokenTransferLogManager.Create(tokenTransferLog)
			if err != nil {
				return fmt.Errorf("create TokenTransferLog error: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
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
	return service.db.Transaction(func(tx *gorm.DB) error {
		workflowManager := do.NewWorkFlowInfoManager(tx)
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

		workflowApproveManager := do.NewWorkFlowApproveManager(tx)
		err = workflowApproveManager.Create(approve)
		if err != nil {
			service.logger.Error("Create WorkFlowApprove error", zap.Error(err))
			return err
		}

		count, err := workflowApproveManager.CountUniqueApprovedAddresses(input.WorkflowID)
		if err != nil {
			service.logger.Error("Count approved addresses error", zap.Error(err))
			return err
		}

		if count >= 2 {
			workflow.Status = do.WorkFlowStatusApproved
			workflow.UpdatedBy = input.ApproverAddr
			workflow.UpdatedAddr = input.ApproverAddr
			workflow.UpdatedTime = time.Now()

			err = workflowManager.Update(workflow)
			if err != nil {
				service.logger.Error("Update workflow status error", zap.Error(err))
				return err
			}
			tokenInfoManager := do2.NewTokenInfoManager(tx)
			tokenInfo, err := tokenInfoManager.GetByID(1)
			if err != nil {
				service.logger.Error("CreateWorkFlow tokenInfoManager GetByID", zap.Error(err))
				return err
			}

			tokenTransferLog := &do2.TokenTransferLog{
				TokenInfoID:     workflow.TokenInfoID,
				WorkflowID:      workflow.ID,
				FromAddress:     "0x0",
				ToAddress:       workflow.ToAddr,
				ContractAddress: tokenInfo.ContractAddress,
				Amount:          0,
				TransferData:    "",
				Status:          do2.StatusPending,
				RetryCount:      0,
				TransactionHash: "",
				CreateBy:        workflow.CreateBy,
				CreateAddr:      workflow.CreateAddr,
				CreatedTime:     time.Now(),
			}

			tokenTransferLogManager := do2.NewTokenTransferLogManager(tx)
			err = tokenTransferLogManager.Create(tokenTransferLog)
			if err != nil {
				return fmt.Errorf("create TokenTransferLog error: %w", err)
			}
		}

		return nil
	})
}
