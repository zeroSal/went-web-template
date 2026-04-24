package app

import (
	"webtemplate/app/bootstrap/module"
	"webtemplate/app/config"

	"go.uber.org/fx"
)

var Container = fx.Module(
	"container",
	fx.Provide(module.IrisProvider),
	fx.Provide(module.AuditLoggerProvider),
	fx.Provide(module.ErrorLoggerProvider),
	fx.Provide(config.LoadEnv),
)