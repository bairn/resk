package resk

import (
	"resk/apis/gorpc"
	_ "resk/apis/web"
	_ "resk/core/envelopes"
	"github.com/bairn/infra"
	"github.com/bairn/infra/base"
	"resk/jobs"
	_ "resk/public/ui"
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
