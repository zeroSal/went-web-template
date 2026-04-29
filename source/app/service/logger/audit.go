package logger

import (
	"context"
	"fmt"
	"webtemplate/app/service/env"

	"github.com/zeroSal/went-logger/logger"
	"go.uber.org/fx"
)

type AuditLogger struct {
	*logger.FileLogger
}

func NewAuditLogger(
	lc fx.Lifecycle,
	env *env.Env,
) *AuditLogger {
	path := fmt.Sprintf("%s/audit.log", env.GetLogsDir())
	l := logger.NewFileLogger(path, "audit", logger.LevelInfo)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return l.Close()
		},
	})

	return &AuditLogger{l}
}