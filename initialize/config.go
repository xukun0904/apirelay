package initialize

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"jhr.com/apirelay/global"
)

func GetEnvInfo(env string) string {
	viper.AutomaticEnv()
	return viper.GetString(env)
}

func InitConfig() {
	zap.S().Info("Start initializing configuration")
	v := viper.New()
	// 配置文件名称及类型
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	// 获取配置文件所在路径
	if apirelayConfigPath := GetEnvInfo("CONFIG_PATH"); apirelayConfigPath != "" {
		zap.S().Info("The path of the configuration file is: ", apirelayConfigPath)
		v.AddConfigPath(apirelayConfigPath)
	} else {
		// 默认在当前目录查找
		v.AddConfigPath(".")
	}
	// 加载配置文件
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panic("Failed to load configuration file: ", err.Error())
	}
	// 序列化成struct
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		zap.S().Panic("Serialization configuration file to struct failed: ", err.Error())
	}
	// 验证配置文件是否正确
	validate := validator.New()
	if err := validate.Struct(&global.ServerConfig); err != nil {
		zap.S().Panic("Configuration verification failed: ", err.Error())
	}
	zap.S().Info("Initial configuration succeeded")
}
