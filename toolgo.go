package toolgo

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lngwu11/toolgo/loggo"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
	"time"
)

type Conf struct {
	ConfigFile ConfigFileConf `mapstructure:"-"`
	Logger     LoggerConf
	Config     interface{}
}

type LoggerConf struct {
	LogLevel         string
	FilePath         string
	FileMaxAge       int // 文件最大保存时间（天）
	FileRotationTime int // 日志切割时间间隔（小时）
}

type ConfigFileConf struct {
	FilePath string
}

var defaultConf = Conf{
	Logger: LoggerConf{
		LogLevel:         "ERROR",
		FilePath:         "logs/toolgo.log",
		FileMaxAge:       30,
		FileRotationTime: 24,
	},
	ConfigFile: ConfigFileConf{
		FilePath: "conf/conf.toml",
	},
}

var cfg = &defaultConf

func (c *Conf) Init() {
	err := initConf(c)
	if err != nil {
		panic(err)
	}
}

func (c *Conf) SetConfigFilePath(path string) *Conf {
	c.ConfigFile.FilePath = path
	return c
}

func (c *Conf) SetLogFilePath(path string) *Conf {
	c.Logger.FilePath = path
	return c
}

func (c *Conf) SetConfig(config interface{}) *Conf {
	c.Config = config
	return c
}

func initLog() (err error) {
	loggo.ResetLogging()
	err = loggo.RegisterWriter(os.Stdout.Name(), loggo.NewSimpleWriter(os.Stdout, nil))
	if err != nil {
		return
	}

	filePath := cfg.Logger.FilePath
	if filePath != "" {
		//获取文件后缀
		fileSuffix := path.Ext(filePath)
		//获取不带后缀的文件名
		filenameOnly := strings.TrimSuffix(filePath, fileSuffix)
		logWriter, _ := rotatelogs.New(
			filenameOnly+".%Y%m%d"+fileSuffix,
			//生成软链 指向最新的日志文件
			rotatelogs.WithLinkName(filePath),
			//文件最大保存时间
			rotatelogs.WithMaxAge(time.Duration(cfg.Logger.FileMaxAge)*24*time.Hour),
			//设置日志切割时间间隔
			rotatelogs.WithRotationTime(time.Duration(cfg.Logger.FileRotationTime)*time.Hour),
		)
		err = loggo.RegisterWriter(path.Base(filePath), loggo.NewSimpleWriter(logWriter, loggo.DefaultFormatter))
		if err != nil {
			return
		}
	}
	return loggo.ConfigureLoggers(cfg.Logger.LogLevel)
}

func initConf(c *Conf) (err error) {
	filePath := c.ConfigFile.FilePath
	if len(filePath) > 0 {
		// 读取配置文件并转化成对应的结构体
		viper.SetConfigFile(filePath)
		if err = viper.ReadInConfig(); err != nil {
			return
		}
		if err = viper.Unmarshal(cfg); err != nil {
			return
		}
		//fmt.Printf("%+v\n", cfg.Config)

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			// 配置文件发生变更之后会调用的回调函数
			_, _ = fmt.Fprintf(os.Stderr, "%s viper.OnConfigChange:%s", time.Now().Format("2006-01-02 15:04:05.000"), e.Name)
		})
	}

	return initLog()
}

func init() {
	_ = initConf(cfg)
}

// Deprecated: use Default instead
func GetDefaultConf() *Conf {
	return cfg
}

func Default() *Conf {
	return cfg
}

func InitLog(level, filePath string) (err error) {
	loggo.ResetLogging()
	err = loggo.RegisterWriter(os.Stdout.Name(), loggo.NewSimpleWriter(os.Stdout, nil))
	if err != nil {
		return
	}
	if filePath != "" {
		var w loggo.Writer
		w, err = loggo.NewLogFileWriter(filePath, 0644, nil)
		if err != nil {
			return
		}
		err = loggo.RegisterWriter(path.Base(filePath), w)
		if err != nil {
			return
		}
	}
	return loggo.ConfigureLoggers(level)
}
