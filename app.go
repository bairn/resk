package resk

import (
	"github.com/bairn/infra"
	"github.com/bairn/infra/base"
	"github.com/bairn/resk/apis/gorpc"
	_ "github.com/bairn/resk/apis/web"
	_ "github.com/bairn/resk/core/envelopes"
	"github.com/bairn/resk/jobs"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	infra.Register(&base.GoRPCStarter{})
	infra.Register(&gorpc.GoRpcApiStarter{})
	infra.Register(&jobs.RefundExpiredJobStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.EurekaStarter{})
	infra.Register(&base.HookStarter{})
}
