package testx

import (
	"github.com/bairn/props/ini"
	"github.com/bairn/props/kvs"
	"github.com/bairn/infra"
	"github.com/bairn/infra/base"
)

func init() {
	file := kvs.GetCurrentFilePath("../brun/config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	app := infra.New(conf)
	app.Start()
}
