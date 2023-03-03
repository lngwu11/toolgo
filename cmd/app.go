package main

import (
	"github.com/lngwu11/toolgo"
	"github.com/lngwu11/toolgo/loggo"
)

var (
	logger = loggo.GetLogger("demo")
)

type DemoConfig struct {
	Name string
	Addr string
	Port int
}

var demoConfig DemoConfig

func main() {
	toolgo.GetDefaultConf().
		SetConfigFilePath("conf/conf.toml").
		SetConfig(&demoConfig).
		Init()

	logger.Debugf("this is a debug message")
	logger.Errorf("this is a error message")
	logger.Debugf("config:%+v", demoConfig)
}
