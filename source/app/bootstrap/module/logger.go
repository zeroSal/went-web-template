package module

import (
	"webtemplate/app/service/logger"

	"go.uber.org/fx"
)

var Loggers = fx.Module(
	"loggers",
	fx.Provide(logger.NewAuditLogger),
	fx.Provide(logger.NewErrorLogger),
)
