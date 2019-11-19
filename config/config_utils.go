package config

import (
	"WeiboSpiderGo/utils"
	"github.com/spf13/viper"
)

var Conf *viper.Viper

func init() {
	Conf = viper.New()

	Conf.SetConfigName("config")

	Conf.AddConfigPath(utils.ExecPath)

	Conf.SetConfigType("yaml")
	if err := Conf.ReadInConfig(); err != nil {
		panic(err)
	}
}
