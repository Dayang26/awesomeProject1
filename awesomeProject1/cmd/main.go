package main

import (
	ginblog "awesomeProject1/internal"
	g "awesomeProject1/internal/global"
	"awesomeProject1/internal/middleware"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func main() {

	configPath := flag.String("c", "resource/config.yaml", "配置文件路径")
	flag.Parse()

	conf := g.ReadConfig(*configPath)

	_ = ginblog.InitLogger(conf)
	db := ginblog.InitDatabase(conf)
	rdb := ginblog.InitRedis(conf)

	gin.SetMode(conf.Server.Mode)

	r := gin.New()
	r.SetTrustedProxies([]string{"*"})

	// 开发模式使用gin自带的日志和恢复中间件，生产模式使用自定义的中间件
	if conf.Server.Mode == "debug" {
		r.Use(gin.Logger(), gin.Recovery())
	} else {
		// todo 生产模式使用自定义的日志和恢复中间件
	}

	r.Use(middleware.CORS())
	r.Use(middleware.WithGormDB(db))
	r.Use(middleware.WithRedisDB(rdb))
	r.Use(middleware.WithCookiesStore(conf.Session.Name, conf.Session.Salt))

	if conf.Upload.OssType == "local" {
		r.Static(conf.Upload.Path, conf.Upload.StorePath)
	}

	serverAddr := conf.Server.Port
	if serverAddr[0] == ':' || strings.HasPrefix(serverAddr, "0.0.0.0:") {
		log.Printf("Serving HTTP on (http://localhost:%s/) ... \n", strings.Split(serverAddr, ":")[1])
	} else {
		log.Printf("Serving HTTP on (http://%s/) ... \n", serverAddr)
	}

	r.Run(serverAddr)
}
