package logger

import (
	"context"
	"fmt"
	"webtemplate/app/service/env"

	"github.com/zeroSal/went-logger/logger"
	"go.uber.org/fx"
)

type ErrorLogger struct {
	*logger.FileLogger
}

func NewErrorLogger(
	lc fx.Lifecycle,
	env *env.Env,
) *ErrorLogger {
	path := fmt.Sprintf("%s/error.log", env.GetLogsDir())
	l := logger.NewFileLogger(path, "error", logger.LevelError)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return l.Close()
		},
	})

	return &ErrorLogger{l}
}
