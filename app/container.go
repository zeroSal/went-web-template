package app

import (
	"template/app/bootstrap/module"
	"template/app/config"

	"go.uber.org/fx"
)

var Container = fx.Module(
	"container",
	fx.Provide(module.IrisProvider),
	fx.Provide(module.AuditLoggerProvider),
	fx.Provide(module.ErrorLoggerProvider),
	fx.Provide(config.LoadEnv),
)