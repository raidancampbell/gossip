package main

type cfg struct {
	Network struct {
		Host     string   `mapstructure:"host"`
		Port     int      `mapstructure:"port"`
		SSL      bool     `mapstructure:"ssl"`
		Channels []string `mapstructure:"channels"`
	} `mapstructure:"network"`
	Nick    string `mapstructure:"nick"`
	Logging struct {
		Level  string `mapstructure:"level"`
		Splunk struct {
			Enabled bool   `mapstructure:"enabled"`
			Host    string `mapstructure:"host"`
			HECPort int    `mapstructure:"hecPort"`
			Token   string `mapstructure:"token"`
		} `mapstructure:"splunk"`
	} `mapstructure:"logging"`
}
