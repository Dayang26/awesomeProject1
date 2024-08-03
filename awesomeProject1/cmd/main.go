package main

import (
	g "awesomeProject1/internal/global"
	"flag"
)

func main() {

	configPath := flag.String("c", "resource/config.yaml", "配置文件路径")
	flag.Parse()

	conf := g.ReadConfig(*configPath)

}
