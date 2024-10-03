package business

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"go-project/business/workflow/dto"
	"go-project/business/workflow/service"
	"go-project/common/types"
	"go-project/common/web"
	"go-project/main/app"
)

func CreateWorkFlow(c *gin.Context) {
	var input dto.WorkflowInfoCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		app.Log.Error("CreateWorkFlow ShouldBindJSON", zap.Any("error", err))
		web.Fail(c, err.Error())
		return
	}

	info, err := service.DefaultService.CreateWorkFlowService(&input)
	if err != nil {
		return
	}

	web.Success(c, info)
}

func WorkFlowList(c *gin.Context) {
	var pageReq types.GenericPageReq[dto.WorkflowInfoCreateDTO]

	pageResp, err := service.DefaultService.PageWorkFlowList(pageReq)
	if err != nil {
		app.Log.Error("WorkFlowList service error", zap.Error(err))
		web.Fail(c, err.Error())
		return
	}

	web.Success(c, pageResp)
}

func WorkFlowApproval(c *gin.Context) {
	var input dto.WorkFlowApprovalDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		app.Log.Error("WorkFlowApproval bind JSON", zap.Error(err))
		web.Fail(c, "input error")
		return
	}

	err := service.DefaultService.ApproveWorkFlow(&input)
	if err != nil {
		app.Log.Error("WorkFlowApproval service error", zap.Error(err))
		web.Fail(c, err.Error())
		return
	}

	web.Success(c, "")
}
