package logic

import (
	"go.uber.org/zap"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"
)

func CreatePost(p *models.Post) (err error) {
	id := snowflake.GenId()
	p.ID = id
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
}

func GetPostDetailById(id int64) (data *models.ApiPostDetail, err error) {

	data = new(models.ApiPostDetail)
	post, err := mysql.GetPostDetailById(id)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById(id) failed",
			zap.Int64("id", id),
			zap.Error(err))
	}
	data.Post = post
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
	}
	data.AuthorName = user.Username
	community, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
	}
	data.CommunityDetail = community

	return data, nil
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById(id) failed",
			zap.Error(err))
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		data = append(data, &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		})
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		data = append(data, &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		})
	}
	return
}

func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetPostCommunityList2(p)
	}

	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}

func GetPostCommunityList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetPostCommunityIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		data = append(data, &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		})
	}
	return
}
