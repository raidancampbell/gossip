package conf

type Cfg struct {
	Network   Network  `mapstructure:"network"`
	Nick      string   `mapstructure:"nick"`
	OwnerNick string   `mapstructure:"ownerNick"`
	Logging   Logging  `mapstructure:"logging"`
	Triggers  Triggers `mapstructure:"triggers"`
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

type Triggers struct {
	Push Push `mapstructure:"push"`
}

type Push struct {
	APIKey       string   `mapstructure:"APIKey"`
	RecipientKey string   `mapstructure:"recipientKey"`
	HighlightOn  []string `mapstructure:"highlightOn"`
}
