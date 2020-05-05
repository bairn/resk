package main

import (
	"github.com/bairn/infra"
	"github.com/bairn/infra/base"
	_ "github.com/bairn/resk"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
}


//func main2() {
//
//	flag.Parse()
//	profile := flag.Arg(0)
//	if profile == "" {
//		profile = "dev"
//	}
//
//	file := kvs.GetCurrentFilePath("boot.ini", 1)
//	log.Info(file)
//
//	conf := ini.NewIniFileCompositeConfigSource(file)
//	if _, err := conf.Get("profile"); err != nil {
//		conf.Set("profile", profile)
//	}
//
//	addr := conf.GetDefault("consul.address", "127.0.0.1:8500")
//	contexts := conf.KeyValue("consul.contexts").Strings()
//
//	log.Info("consul address:", addr)
//	log.Info("consul contexts:", contexts)
//
//	consulConf := consul.NewCompositeConsulConfigSourceByType(contexts, addr, kvs.ContentIni)
//	consulConf.Add(conf)
//
//
//	base.InitLog(kvs.ConfigSource(consulConf))
//	app := infra.New(consulConf)
//	app.Start()
//}
