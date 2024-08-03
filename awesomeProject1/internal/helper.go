package ginblog

import (
	g "awesomeProject1/internal/global"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"log/slog"
	"os"
	"time"
)

// InitLogger 初始化日志系统。
// 它根据配置文件中的日志级别和格式设置来初始化一个日志记录器。
// 参数conf是全局配置结构体的指针，用于获取日志相关的配置。
// 返回值是初始化后的日志记录器实例。
func InitLogger(conf *g.Config) *slog.Logger {
	// 初始化日志级别变量，根据配置文件中的日志级别来设置。
	var level slog.Level
	switch conf.Log.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// 配置日志处理器的选项，包括是否添加源代码信息、日志级别，以及自定义的时间格式化函数。
	option := &slog.HandlerOptions{
		AddSource: false,
		Level:     level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime))
				}
			}
			return a
		},
	}

	// 根据配置文件中的日志格式选择合适的日志处理器。
	var handle slog.Handler
	switch conf.Log.Format {
	case "json":
		handle = slog.NewJSONHandler(os.Stdout, option)
	case "text":
		// 如果配置为"text"或者未指定，默认使用文本格式化处理器。
		fallthrough
	default:
		handle = slog.NewTextHandler(os.Stdout, option)
	}

	// 使用选定的处理器创建日志记录器实例，并设置为默认日志记录器。
	logger := slog.New(handle)
	slog.SetDefault(logger)
	return logger
}

func initDatabase(conf *g.Config) *gorm.DB {
	dbType := conf.DbType()
	dsn := conf.DbDSN()

	var db *gorm.DB
	var err error

	var level logger.LogLevel
	switch conf.Server.DbLogMod {
	case "silent":
		level = logger.Silent
	case "info":
		level = logger.Info
	case "warn":
		level = logger.Warn
	case "error":
		fallthrough
	default:
		level = logger.Error
	}

	config := &gorm.Config{
		Logger:                                   logger.Default.LogMode(level),
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	switch dbType {
	case "mysql":
		gorm.Open(mysql.Open(dsn), config)
	case "sqlite":
	default:
		log.Fatal("不支持的数据库类型: ", dbType)
	}

	return nil
}
