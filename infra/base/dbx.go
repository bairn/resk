package base

import (
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/tietang/props/kvs"
	"resk/infra"
	"resk/infra/logrus"
)

var database *dbx.Database

func DbxDatabase() *dbx.Database {
	Check(database)
	return database
}

type DbxDatabaseStarter struct {
	infra.BaseStarter
}

func (s *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {
	conf := ctx.Props()

	settings := dbx.Settings{
		//LoggingEnabled:true,
	}
	err := kvs.Unmarshal(conf, &settings, "mysql")
	if err != nil {
		panic(err)
	}

	log.Info("mysql.conn url:", settings.ShortDataSourceName())
	dbx , err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	//log.Info(dbx.Ping())
	dbx.SetLogger(logrus.NewUpperLogrusLogger())

	database = dbx
}