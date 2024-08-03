package main

import (
	ginblog "awesomeProject1/internal"
	g "awesomeProject1/internal/global"
	"flag"
)

func main() {

	configPath := flag.String("c", "resource/config.yaml", "配置文件路径")
	flag.Parse()

	conf := g.ReadConfig(*configPath)

	_ = ginblog.InitLogger(conf)
	db := ginblog.InitDatabase(conf)
	rdb := ginblog.InitRedis(conf)

}
