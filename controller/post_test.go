package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()
	url := "/api/v1/post"
	engine.POST(url, CreatePostHandler)

	body := `{
    "community_id": 3,
    "title": "3、这是我刚刚加入的第er天",
    "content": "xx5555533xxx" 
}`
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, request)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "请登录")

	data := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), data); err != nil {
		t.Fatalf("json.Unmarshal w.body() failed err:%v\n", err)
	}
	assert.Equal(t, data.Code, CodeNeedLogin)
}
