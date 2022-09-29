package logic

import (
	"database/sql"
	"errors"
	"web_app/dao/mysql"
	"web_app/models"
	"web_app/pkg/jwt"
	"web_app/pkg/snowflake"
)

var (
	ErrorUserNotExists   = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
)

func SignUp(p *models.ParamSignUp) (err error) {

	if err := mysql.CheckUserExists(p.UserName); err != nil {
		return err
	}

	userId := snowflake.GenId()
	user := &models.User{
		UserID:   userId,
		Username: p.UserName,
		Password: p.Password,
	}
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (token string, err error) {

	var user models.User
	user, err = mysql.FindUser(p.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrorUserNotExists
		}
		return "", err
	}

	if user.Password != mysql.EncryptPassword(p.Password) {
		return "", ErrorInvalidPassword
	}

	return jwt.GetToken(user.UserID, user.Username)
}
