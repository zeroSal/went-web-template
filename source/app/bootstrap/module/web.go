package module

import (
	"webtemplate/app/registry"

	"github.com/zeroSal/went-web/config"
	"github.com/zeroSal/went-web/factory"
	"go.uber.org/fx"
)

var irisConfig = config.NewIris(
	true,
	registry.Controller,
)

var Web = fx.Module(
	"web",
	fx.Supply(irisConfig),
	fx.Provide(factory.SecurityFactory),
	fx.Provide(factory.IrisFactory),
);