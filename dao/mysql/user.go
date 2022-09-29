package mysql

import (
	"crypto/md5"
	"encoding/hex"
	"web_app/models"
)

// CheckUserExists 根据username检查用户是否存在
func CheckUserExists(username string) (err error) {
	strSql := `select count(user_id) from user where username = ?`
	var count int
	if err = db.Get(&count, strSql, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExists
	}
	return
}

// InsertUser 插入一条新的用户
func InsertUser(user *models.User) (err error) {
	strSql := `insert into user(user_id,username,password) values(?,?,?)`
	user.Password = EncryptPassword(user.Password)
	_, err = db.Exec(strSql, user.UserID, user.Username, user.Password)
	return
}

func FindUser(username string) (user models.User, err error) {
	strSql := `select user_id, username, password from user where username = ?`
	err = db.Get(&user, strSql, username)
	return
}

const secret = "sunx@111"

func EncryptPassword(txt string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(txt)))
}
