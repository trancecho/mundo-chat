package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func ConfigInit() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("Error reading config file:", err)
		//检查配置文件是否存在
		if _, statErr := os.Stat("./config/config.yaml"); os.IsNotExist(statErr) {
			log.Println("./config/config.yaml not found")
		}
	}
}
