package module

import (
	"webtemplate/app/service/provider"

	"go.uber.org/fx"
)

var Provider = fx.Module(
	"provider",
	fx.Provide(provider.NewSession),
)