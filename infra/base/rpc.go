package base

import (
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"reflect"
	"resk/infra"
)

var rpcServer *rpc.Server

func RpcServer() *rpc.Server {
	Check(rpcServer)
	return rpcServer
}

func RpcRegister(ri interface{}) {
	typ := reflect.TypeOf(ri)
	log.Infof("goRPC Register: %s", typ.String())
	RpcServer().Register(ri)
}

type GoRPCStarter struct {
	infra.BaseStarter
	server *rpc.Server
}

func (s *GoRPCStarter) Init(ctx infra.StarterContext) {
	s.server = rpc.NewServer()
	rpcServer = s.server
}

func (s *GoRPCStarter) Start(ctx infra.StarterContext) {
	port := ctx.Props().GetDefault("app.rpc.port", "8082")
	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Panic(err)
	}
	log.Info("tcp port listened form rpf:", port)
	go s.server.Accept(listener)
}
