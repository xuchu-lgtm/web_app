package mysql

import (
	"database/sql"
	"go.uber.org/zap"
	"web_app/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	strSql := `select community_id, community_name from community limit 20`
	if err := db.Select(&communityList, strSql); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailById(id int64) (community *models.CommunityDetail, err error) {
	communityDetail := new(models.CommunityDetail)
	strSql := `select community_id,community_name,introduction,create_time from community where community_id = ?`
	if err := db.Get(communityDetail, strSql, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrorInvalidId
		}
	}
	return communityDetail, err
}
