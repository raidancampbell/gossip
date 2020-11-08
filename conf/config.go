package conf

type Cfg struct {
	Network Network `mapstructure:"network"`
	Nick    string  `mapstructure:"nick"`
	Logging Logging `mapstructure:"logging"`
}

type Network struct {
	Host     string   `mapstructure:"host"`
	Port     int      `mapstructure:"port"`
	SSL      bool     `mapstructure:"ssl"`
	Channels []string `mapstructure:"channels"`
}

type Logging struct {
	Level  string        `mapstructure:"level"`
	Splunk SplunkLogging `mapstructure:"splunk"`
}

type SplunkLogging struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	HECPort int    `mapstructure:"hecPort"`
	Token   string `mapstructure:"token"`
}