package business

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-project/chain/eth"
	"go-project/main/log"
)

type Route struct {
	DB          *gorm.DB
	Log         *log.ZapLogger
	ERC20Client eth.TestErc20Client
}

func (r *Route) Register(engine *gin.Engine) {
	root := engine.Group("")

	root.POST("/workflow/create", func(c *gin.Context) {
		CreateWorkFlow(c, r.DB, r.Log, r.ERC20Client)
	})
	root.GET("/workflow/page", func(c *gin.Context) {
		WorkFlowList(c, r.DB, r.Log)
	})
	root.POST("/workflow/approve", func(c *gin.Context) {
		WorkFlowApproval(c, r.DB, r.Log)
	})
}
