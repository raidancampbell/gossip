package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	conf = &cfg{}
)

func initialize() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.Unmarshal(conf)
}

func main() {
	initialize()
	logrus.Info(conf.Nick)
	fmt.Println(conf.Network.Host)
}
