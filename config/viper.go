package config

import (
	"github.com/spf13/viper"
	"log"
)

func ConfigInit() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Fatal error config file: ./config/app.yaml")
	}
}
