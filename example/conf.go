package main

import (
	"fmt"
	"github.com/bairn/props/ini"
	"github.com/bairn/props/kvs"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileConfigSource(file)
	port := conf.GetIntDefault("app.server.port", 18080)
	fmt.Println(port)
}
