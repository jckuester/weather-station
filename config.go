package main

import (
	"github.com/spf13/viper"
	"log"
)

var vip = viper.New()

func initConfig() {
	vip.SetConfigType("yaml")
	vip.SetConfigName("weather-station")
	vip.AddConfigPath("/etc/weather-station/")
	vip.AddConfigPath("$HOME/.config/weather-station")
	vip.AddConfigPath(".")
	vip.SetDefault("sensors", map[string]string{})
}

func readConfig() {
	err := vip.ReadInConfig()
	if err != nil {
		log.Println(err, "Error reading config file. Running in scanning mode.")
	}
}

func loadConfig() {
	vip = viper.New()
	initConfig()
	readConfig()
}
