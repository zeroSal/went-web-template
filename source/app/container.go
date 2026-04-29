package app

import (
	"webtemplate/app/bootstrap/module"
	"webtemplate/app/service/env"

	"go.uber.org/fx"
)

var Container = fx.Module(
	"container",
	module.Provider,
	module.Web,
	module.Loggers,
	fx.Provide(env.LoadEnv),
)
