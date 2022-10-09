package main

import (
	"context"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/pkg/consul"
	"web_app/pkg/jaeger"
	"web_app/pkg/snowflake"
	"web_app/routes"
	"web_app/settings"
)

// @title 标题2
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 这里写接口服务的host
// @BasePath 这里写base path
func main() {
	var configFile string
	flag.StringVar(&configFile, "configFile", "./conf/config.yaml", "配置文件")
	flag.Parse()

	// 加载配置
	if err := settings.Init(configFile); err != nil {
		fmt.Printf("init settings failed, err: %v\n", err)
		return
	}
	// 初始化日志
	defer zap.L().Sync() //把缓冲区的日志追加到里面
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}

	// 注册consul
	consul, err := consul.NewConsul(settings.Conf.ConsulConfig)
	if err != nil {
		fmt.Printf("init consul failed, err: %v\n", err)
	}
	if err := consul.Init(settings.Conf); err != nil {
		fmt.Printf("init consul failed, err: %v\n", err)
		return
	}

	// 初始化MySql连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err: %v\n", err)
		return
	}
	defer mysql.Close()

	// 初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err: %v\n", err)
		return
	}
	defer redis.Close()

	// 初始化jaeger
	if err := jaeger.Init(settings.Conf.Name, settings.Conf.JaegerConfig); err != nil {
		fmt.Printf("init jaeger failed, err: %v\n", err)
		return
	}
	defer jaeger.Close()

	//初始化snowflake
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineId); err != nil {
		fmt.Printf("init snowflake failed, err: %v\n", err)
		return
	}

	// 5. 注册路由
	router := routes.SetupRouter(settings.Conf.Mode)
	// 6. 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 退出时注销服务
	consul.Deregister(settings.Conf.UUID)
	zap.L().Info("Shutdown Server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
