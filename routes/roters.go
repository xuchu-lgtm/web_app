package routes

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
	"web_app/controller"
	_ "web_app/docs"
	"web_app/logger"
	"web_app/middlewares"
)

func SetupRouter(mode string) *gin.Engine {

	// gin 设置发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	//初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
	}

	r := gin.New()
	// 告诉gin框架模板文件引用的静态文件去哪里找
	r.Static("/static", "static")
	// 告诉gin框架去哪里找模板文件
	r.LoadHTMLGlob("templates/*")

	// 令牌桶=> 每两秒钟可以取一个令牌
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 10000))

	r.GET("/", controller.IndexHandler)
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	{
		//注册
		v1.POST("/signup", controller.SignUpHandler)
		//登录
		v1.POST("/login", controller.LoginHandler)

		//注册中间件
		v1.Use(middlewares.JWTAuthMiddleware())
		{
			v1.GET("/community", controller.CommunityHandler)
			v1.GET("/community/:id", controller.CommunityDetailHandler)

			v1.POST("/post", controller.CreatePostHandler)
			v1.GET("/post/:id", controller.GetPostDetailHandler)
			v1.GET("/posts", controller.GetPostListHandler)
			v1.GET("/posts2", controller.GetPostListHandler2)

			v1.POST("/vote", controller.PostVoteController)
		}
	}

	pprof.Register(r) //注册pprof相关路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
