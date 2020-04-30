package gorpc

import (
	"resk/infra"
	"resk/infra/base"
)

type GoRpcApiStarter struct {
	infra.BaseStarter
}

func (g *GoRpcApiStarter) Init(ctx infra.StarterContext) {
	base.RpcRegister(new(EnvelopeRpc))
}
