package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init() (err error) {
	viper.SetConfigName("config") //指定配置文件名称（不需要带后缀）
	viper.SetConfigType("yaml")   //指定配置文件类型
	viper.AddConfigPath(".")      //指定查找配置文件的路径（这里使用相对路径）
	err = viper.ReadInConfig()    //读取配置信息
	if err != nil {
		zap.L().Error("viper.ReadInConfig() failed, err:%v\n", zap.Error(err))
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info("配置文件修改了... %s", zap.String("name", in.Name))
	})
	return
}
