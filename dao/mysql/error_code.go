package mysql

import "errors"

var (
	ErrorUserExists = errors.New("用户已存在")
	ErrorInvalidId  = errors.New("无效的Id")
)
