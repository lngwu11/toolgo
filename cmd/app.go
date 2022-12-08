package main

import (
	"github.com/lngwu11/toolgo"
)

var (
	logger = toolgo.GetLogger("demo")
)

type DemoConfig struct {
	Name string
	Addr string
	Port int
}

var demoConfig DemoConfig

func main() {
	cfg := toolgo.GetDefaultConf()
	cfg.SetConfig(&demoConfig)
	_ = cfg.Init()

	logger.Debugf("this is a debug message")
	logger.Errorf("this is a error message")
	logger.Debugf("config:%+v", demoConfig)
}
