package models

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 定义注册请求参数结构体
type ParamSignUp struct {
	UserName   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 定义登录请求参数结构体
type ParamLogin struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VoteData struct {
	PostID    string `json:"post_id" binding:"required"`
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` //oneof只能是1,0,-1
}

type ParamPostList struct {
	Page        int64  `json:"page" form:"page"`                 // 页码
	Size        int64  `json:"size" form:"size"`                 // 每页数据量
	Order       string `json:"order" form:"order"`               // 排序依据
	CommunityID int64  `json:"community_id" form:"community_id"` // 可以为空
}
