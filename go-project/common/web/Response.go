package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		http.StatusOK,
		data,
		"ok",
	})
}

func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		http.StatusInternalServerError,
		nil,
		msg,
	})
}

func FailV2(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		code,
		nil,
		msg,
	})
}
