package business

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-project/business/workflow/dto"
	"go-project/business/workflow/service"
	"go-project/common/types"
	"go-project/common/web"
	"go-project/main/log"
)

func CreateWorkFlow(c *gin.Context, db *gorm.DB, log *log.ZapLogger) {
	var input dto.WorkflowInfoCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error("CreateWorkFlow ShouldBindJSON", zap.Any("error", err))
		web.Fail(c, err.Error())
		return
	}

	info, err := service.NewService(log, db).CreateWorkFlowService(&input)
	if err != nil {
		return
	}

	web.Success(c, info)
}

func WorkFlowList(c *gin.Context, db *gorm.DB, log *log.ZapLogger) {
	var pageReq types.GenericPageReq[dto.WorkflowInfoCreateDTO]

	pageResp, err := service.NewService(log, db).PageWorkFlowList(pageReq)
	if err != nil {
		log.Error("WorkFlowList service error", zap.Error(err))
		web.Fail(c, err.Error())
		return
	}

	web.Success(c, pageResp)
}

func WorkFlowApproval(c *gin.Context, db *gorm.DB, log *log.ZapLogger) {
	var input dto.WorkFlowApprovalDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error("WorkFlowApproval bind JSON", zap.Error(err))
		web.Fail(c, "input error")
		return
	}

	err := service.NewService(log, db).ApproveWorkFlow(&input)
	if err != nil {
		log.Error("WorkFlowApproval service error", zap.Error(err))
		web.Fail(c, err.Error())
		return
	}

	web.Success(c, "")
}
