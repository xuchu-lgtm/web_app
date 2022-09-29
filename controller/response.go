package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
 "code": 10001,
 "msg":  xxx
 "data": {}
}
*/

func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  code.Msg(),
		"data": nil,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": CodeSuccess,
		"msg":  CodeSuccess.Msg(),
		"data": data,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}
