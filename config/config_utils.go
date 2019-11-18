package config

import (
"github.com/spf13/viper"
)

var Conf *viper.Viper

func init() {
	Conf = viper.New()

	Conf.SetConfigName("config")

	Conf.AddConfigPath("./")

	Conf.SetConfigType("yaml")
	if err := Conf.ReadInConfig(); err != nil {
		panic(err)
	}
}
