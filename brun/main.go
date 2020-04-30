package main

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	_ "resk"
	"resk/infra"
	"resk/infra/base"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
}
