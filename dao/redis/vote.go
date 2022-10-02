package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"math"
	"strconv"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postId, communityID int64) error {

	pipeline := client.TxPipeline()

	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})

	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScore), redis.Z{
		Score:  0,
		Member: postId,
	})

	pipeline.SAdd(getRedisKey(KeyCommunitySetPf+strconv.Itoa(int(communityID))), postId)

	_, err := pipeline.Exec()

	return err
}

func VoteForPost(userID, postID string, value float64) error {

	//去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTime), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	// 先查当前用户给当前帖子的投票记录
	ov := client.ZScore(getRedisKey(KeyPostVotedPf+postID), userID).Val()
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) // 计算两次投票的差值

	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScore), op*diff*scorePerVote, postID)

	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedPf+postID), postID)
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedPf+postID), redis.Z{
			Score:  value, //赞成票还是反对票
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
