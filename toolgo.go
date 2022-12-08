package toolgo

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
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
	IsWriteFile      bool
	FileMaxAge       int // 文件最大保存时间（天）
	FileRotationTime int // 日志切割时间间隔（小时）
	FilePath         string
}

type ConfigFileConf struct {
	FilePath string
}

var DefaultConf = Conf{
	Logger: LoggerConf{
		LogLevel:         "WARNING",
		IsWriteFile:      true,
		FileMaxAge:       7,
		FileRotationTime: 24,
		FilePath:         "logs/logfile.out",
	},
	ConfigFile: ConfigFileConf{
		FilePath: "conf/conf.toml",
	},
}

var cfg = &DefaultConf

func (c *Conf) Init() (err error) {
	cfg = c
	//fmt.Printf("%+v\n", cfg)

	filePath := c.ConfigFile.FilePath
	if len(filePath) > 0 {
		// 读取配置文件并转化成对应的结构体
		viper.SetConfigFile(filePath)
		err = viper.ReadInConfig()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "viper.ReadInConfig:", err)
			return
		}
		if err = viper.Unmarshal(cfg); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "viper.Unmarshal:", err)
			return
		}
		//fmt.Printf("%+v\n", cfg.Config)

		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			// 配置文件发生变更之后会调用的回调函数
			_, _ = fmt.Fprintln(os.Stderr, "viper.OnConfigChange:", e.Name)
		})
	}

	initLog()

	return nil
}

func (c *Conf) SetConfigFilePath(path string) {
	c.ConfigFile.FilePath = path
}

func (c *Conf) SetLogFilePath(path string) {
	c.Logger.FilePath = path
}

func (c *Conf) SetConfig(config interface{}) {
	c.Config = config
}

func initLog() {
	ResetLogging()
	_ = RegisterWriter(os.Stdout.Name(), NewSimpleWriter(os.Stdout, nil))
	if cfg.Logger.IsWriteFile {
		filePath := cfg.Logger.FilePath
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
		_ = RegisterWriter(path.Base(filePath), NewSimpleWriter(logWriter, DefaultFormatter))
	}
	_ = ConfigureLoggers(cfg.Logger.LogLevel)
}

func init() {
	_ = cfg.Init()
}

func GetDefaultConf() *Conf {
	return cfg
}
