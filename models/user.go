package models

type User struct {
	// `db:"user_id,string"` 其中string可以解决前端int64值过大导致失真的问题，如果前端传输string类型的时候也是支持反序列化
	UserID   int64  `db:"user_id,string"`
	Username string `db:"username"`
	Password string `db:"password"`
}
