package business

import "github.com/gin-gonic/gin"

func SetApiGroupRoutes(router *gin.RouterGroup) {
	router.POST("/workflow/create", CreateWorkFlow)
	router.GET("/workflow/page", WorkFlowList)
	router.POST("/workflow/approve", WorkFlowApproval)
}
