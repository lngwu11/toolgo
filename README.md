# toolgo

可以直接集成到项目中用于读取配置文件、打印日志等，方便快捷。

配置文件可以为toml、yaml等其他格式。默认包含Logger日志模块，也可以自己重新定义。格式如下：

```toml
[Logger]
LogLevel = "debug"
FilePath = "logs/demo.log"

[Config]
Name = "Demo"
Addr = "127.0.0.1"
Port = 8888
```

使用时传递配置文件路径和配置结构，然后调用初始化。示例代码：

```golang
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
```

