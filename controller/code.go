package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExists
	CodeUserNotExists
	CodeInvalidPassword
	CodeServerBusy
	CodeRegisterError
	CodeNeedLogin
	CodeInvalidToken
	CodeInvalidId
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数无效",
	CodeUserExists:      "用户名已存在",
	CodeUserNotExists:   "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",
	CodeRegisterError:   "注册失败",
	CodeNeedLogin:       "请登录",
	CodeInvalidToken:    "无效的Token",
	CodeInvalidId:       "无效的Id",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
