package redis

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"web_app/models"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	// 查询索引的起点
	start := (page - 1) * size
	end := start + size - 1

	return client.ZRevRange(key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {

	key := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScore)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

func GetPostCommunityIDsInOrder(p *models.ParamPostList) ([]string, error) {

	orderKey := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScore)
	}

	cKey := getRedisKey(KeyCommunitySetPf + strconv.Itoa(int(p.CommunityID)))
	key := p.Order + strconv.Itoa(int(p.CommunityID))

	if client.Exists(key).Val() < 1 {

		pipeline := client.Pipeline()

		pipeline.ZInterStore(key,
			redis.ZStore{
				Aggregate: "MAX",
			},
			cKey, orderKey)
		pipeline.Expire(key, 60*time.Second)

		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

func GetPostVoteData(ids []string) (data []int64, err error) {
	/*data = make([]int64, 0, len(ids))
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPf + id)
		//查找key中分数
		v := client.ZCount(key, "1", "1").Val()
		data = append(data, v)
	}*/
	// 使用pipeline一次发送数据，减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPf + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
