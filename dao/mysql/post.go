package mysql

import (
	"github.com/jmoiron/sqlx"
	"strings"
	"web_app/models"
)

func CreatePost(p *models.Post) (err error) {
	strSql := `insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)`
	_, err = db.Exec(strSql, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostDetailById(id int64) (data *models.Post, err error) {
	data = new(models.Post)
	strSql := `select post_id, author_id,community_id,status,title,content,create_time from post where post_id = ?`
	err = db.Get(data, strSql, id)
	return
}

func GetPostList(page, size int64) (data []*models.Post, err error) {
	strSql := `select post_id, author_id,community_id,status,title,content,create_time from post order by id desc limit ?,?`
	data = make([]*models.Post, 0, 2)
	err = db.Select(&data, strSql, (page-1)*size, size)
	return
}

func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	postList = make([]*models.Post, 0, len(ids))
	strSql := `select post_id, title, content, author_id, community_id, create_time 
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)
	`
	//http://www.liwenzhou.com/posts/sqlx/
	query, args, err := sqlx.In(strSql, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}

	query = db.Rebind(query)
	err = db.Select(&postList, query, args...) //!!!!!!!
	return
}
