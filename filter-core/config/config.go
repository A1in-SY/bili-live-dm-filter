package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

var Conf *Config

func init() {
	Conf = &Config{
		LogConf: &LogConfig{
			Level:      "debug",
			FileName:   "filter-conn.log",
			MaxSize:    200,
			MaxAge:     0,
			MaxBackups: 0,
			IsStdout:   true,
		},
		ConnConf: &ConnConfig{
			ForceAuth:         true,
			AuthUid:           0,
			AuthCookie:        "",
			HeartbeatInterval: 30 * time.Second,
		},
	}
	var confPath string
	flag.StringVar(&confPath, "conf", "../config.toml", "config path, eg: -conf /path/to/your/config.toml")
	flag.Parse()

	if _, err := os.Stat(confPath); err != nil {
		panic(err)
	}
	data, err := os.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(data, Conf)
	if err != nil {
		panic(err)
	}
}

type Config struct {
	LogConf  *LogConfig
	ConnConf *ConnConfig
}

type LogConfig struct {
	Level    string
	FileName string
	// 单位MB
	MaxSize    int
	MaxAge     int
	MaxBackups int
	IsStdout   bool
}

type ConnConfig struct {
	ForceAuth         bool
	AuthUid           int64
	AuthCookie        string
	HeartbeatInterval time.Duration
}
