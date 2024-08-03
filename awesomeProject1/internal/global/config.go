package g

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	Server struct {
		Mode          string
		Port          string
		DbType        string
		DbAutoMigrate bool
		DbLogMode     string
	}

	Log struct {
		Level     string
		Prefix    string
		Format    string
		Directory string
	}

	JWT struct {
		Secret string
		Expire int64
		Issuer string
	}

	Mysql struct {
		Host     string
		Port     string
		Config   string
		Dbname   string
		Username string
		Password string
	}

	SQLite struct {
		Dsn string // Data Source Name
	}
	Redis struct {
		DB       int    // 指定 Redis 数据库
		Addr     string // 服务器地址:端口
		Password string // 密码
	}
	Session struct {
		Name   string
		Salt   string
		MaxAge int
	}
	Email struct {
		To       string // 收件人 多个以英文逗号分隔 例：a@qq.com,b@qq.com
		From     string // 发件人 要发邮件的邮箱
		Host     string // 服务器地址, 例如 smtp.qq.com 前往要发邮件的邮箱查看其 smtp 协议
		Secret   string // 密钥, 不是邮箱登录密码, 是开启 smtp 服务后获取的一串验证码
		Nickname string // 发件人昵称, 通常为自己的邮箱名
		Port     int    // 前往要发邮件的邮箱查看其 smtp 协议端口, 大多为 465
		IsSSL    bool   // 是否开启 SSL
	}
	Captcha struct {
		SendEmail  bool // 是否通过邮箱发送验证码
		ExpireTime int  // 过期时间
	}
	Upload struct {
		// Size      int    // 文件上传的最大值
		OssType   string // local | qiniu
		Path      string // 本地文件访问路径
		StorePath string // 本地文件存储路径
	}

	Qiniu struct {
		ImgPath       string // 外链链接
		Zone          string // 存储区域
		Bucket        string // 空间名称
		AccessKey     string // 秘钥AK
		SecretKey     string // 秘钥SK
		UseHTTPS      bool   // 是否使用https
		UseCdnDomains bool   // 上传是否使用 CDN 上传加速
	}
}

var Conf *Config

func GetConfig() *Config {
	if Conf == nil {
		log.Panicln("配置文件未初始化")
		return nil
	}
	return Conf
}

func ReadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		panic("配置文件读取失败:" + err.Error())
	}

	if err := v.Unmarshal(&Conf); err != nil {
		panic("配置文件反序列化失败: " + err.Error())
	}

	log.Println("配置文件加载成功:", path)
	return Conf
}

func (*Config) DbType() string {
	if Conf.Server.DbType == "" {
		Conf.Server.DbType = "sqlite"
	}
	return Conf.Server.DbType
}

// 数据库连接字符串
func (*Config) DbDSN() string {
	switch Conf.Server.DbType {
	case "mysql":
		conf := Conf.Mysql
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?%s",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname, conf.Config,
		)
	case "sqlite":
		return Conf.SQLite.Dsn
	// 默认使用 sqlite, 并且使用内存数据库
	default:
		Conf.Server.DbType = "sqlite"
		if Conf.SQLite.Dsn == "" {
			Conf.SQLite.Dsn = "file::memory:"
		}
		return Conf.SQLite.Dsn
	}
}
