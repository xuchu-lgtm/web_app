package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/logic"
	"web_app/models"
)

func SignUpHandler(c *gin.Context) {
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("请求参数有误", zap.Error(err))

		// 判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	if err := logic.SignUp(p); err != nil {
		if errors.Is(err, mysql.ErrorUserExists) {
			ResponseError(c, CodeUserExists)
			return
		}
		ResponseError(c, CodeRegisterError)
		return
	}
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("请求参数有误", zap.Error(err))

		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error(" logic.Login failed err,", zap.String("username", p.UserName), zap.Error(err))

		if errors.Is(err, logic.ErrorUserNotExists) {
			ResponseError(c, CodeUserNotExists)
			return
		}
		if errors.Is(err, logic.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		return
	}

	ResponseSuccess(c, token)
}
