package business

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go-project/business/workflow/dto"
	"go-project/business/workflow/service"
	"go-project/chain/eth"
	globalconst "go-project/common"
	"go-project/common/types"
	"go-project/common/web"
	"go-project/main/log"
)

func CreateWorkFlow(c *gin.Context, db *gorm.DB, log *log.ZapLogger, ERC20Client eth.TestErc20Client) {
	var input dto.WorkflowInfoCreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error("CreateWorkFlow ShouldBindJSON", zap.Any("error", err))
		web.Fail(c, err.Error())
		return
	}

	privateKey, err := crypto.HexToECDSA(globalconst.OWNER_PRV_KEY)
	if err != nil {
		log.Error("Failed to OWNER_PRV_KEY HexToECDSA", zap.Error(err))
		web.Fail(c, "Failed to OWNER_PRV_KEY HexToECDSA")
		return
	}
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	balance, err := ERC20Client.BalanceOf(c.Request.Context(), fromAddress)
	if err != nil {
		log.Error("Failed to get balance", zap.Error(err))
		web.Fail(c, "Failed to get balance")
		return
	}

	requiredBalance := big.NewInt(2e6)
	if balance.Cmp(requiredBalance) < 0 {
		log.Error("Insufficient balance", zap.String("address", fromAddress.Hex()), zap.String("balance", balance.String()))
		web.Fail(c, "Insufficient balance")
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
