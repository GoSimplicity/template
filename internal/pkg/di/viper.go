package di

import (
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitViper() error {
	configFile := pflag.String("config", "", "配置文件路径")
	pflag.Parse()

	// 优先使用命令行参数指定的配置文件
	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	} else {
		// 检查环境变量 ENV
		env := os.Getenv("ENV")
		switch env {
		case "online":
			viper.SetConfigFile("config/config.production.yaml")
		case "dev":
			viper.SetConfigFile("config/config.deployment.yaml")
		default:
			viper.SetConfigFile("config/config.deployment.yaml")
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
