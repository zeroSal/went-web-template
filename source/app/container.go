package app

import (
	"webtemplate/app/bootstrap"
	"webtemplate/app/service/env"
	"webtemplate/app/service/logger"
	"webtemplate/registry"

	wentweb "github.com/zeroSal/went-web"
	"github.com/zeroSal/went-web/controller"

	"go.uber.org/fx"
)

var Container fx.Option

var services = []fx.Option{
	fx.Provide(logger.NewAuditLogger),
	fx.Provide(logger.NewErrorLogger),
	fx.Provide(env.Load),
	wentweb.Bundle,
	bootstrap.Init,
}

func init() {
	opts := services
	for _, constructor := range registry.Controller.All() {
		opts = append(opts, fx.Provide(
			fx.Annotate(
				constructor,
				fx.As(new(controller.Interface)),
				fx.ResultTags(`group:"controllers"`),
			),
		))
	}
	Container = fx.Options(opts...)
}
