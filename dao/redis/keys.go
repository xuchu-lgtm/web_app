package redis

const (
	KeyPrefix         = "webapp:"
	KeyPostTime       = "post:time"
	KeyPostScore      = "post:score"
	KeyPostVotedPf    = "post:voted:"
	KeyCommunitySetPf = "community:" //set保存每个分区下的帖子
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
