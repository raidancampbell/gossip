package main

import (
	"fmt"
	"github.com/raidancampbell/gossip/conf"
	"github.com/raidancampbell/gossip/gossip"
	"github.com/raidancampbell/gossip/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var (
	cfg = &conf.Cfg{}
)

func initialize() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.Unmarshal(cfg)

	// http://splunk.autok8s.raidancampbell.com/en-US/app/search/search
	if cfg.Logging.Splunk.Enabled {
		logrus.AddHook(
			logging.NewSplunkHook(
				http.DefaultClient,
				fmt.Sprintf("http://%s:%d/services/collector", cfg.Logging.Splunk.Host, cfg.Logging.Splunk.HECPort),
				cfg.Logging.Splunk.Token,
				1*time.Second,
				100))
	}

	lvl, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	logrus.SetLevel(lvl)

}

func main() {
	fmt.Println("initializing...")
	initialize()
	logrus.Info("initialized")
	logrus.Info(cfg.Nick)
	fmt.Println(cfg.Network.Host)
	g := gossip.New(cfg)

	// blocking
	g.Begin()
}
