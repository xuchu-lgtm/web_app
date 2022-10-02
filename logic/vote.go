package logic

import (
	"go.uber.org/zap"
	"strconv"
	"web_app/dao/redis"
	"web_app/models"
)

func VoteForPost(userId int64, p *models.VoteData) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userId", userId),
		zap.String("postId", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userId)), p.PostID, float64(p.Direction))
}
