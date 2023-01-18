package config

import (
	"github.com/spf13/viper"
	"schedule/util/log"
)

var Viper *viper.Viper

func InitConfig(path string) error {
	Viper = viper.New()
	if path == "" {
		Viper.SetConfigFile("../../conf/app.toml")
	} else {
		Viper.SetConfigFile(path)
	}
	if err := Viper.ReadInConfig(); err != nil {
		log.ErrLogger.Printf("init config error:%+v", err)
		return err
	}
	return nil
}
