package di

import (
	"github.com/zeroSal/go-semantic-log/logger"
	"go.uber.org/fx"
)

var Container = fx.Module(
	"container",
	fx.Provide(logger.NewConsoleLogger),
)