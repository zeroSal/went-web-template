package bootstrap

import (
	"fmt"
	"os"
	"webtemplate/app/service/env"
	"webtemplate/app/service/logger"

	"go.uber.org/fx"
)

var Init = fx.Options(
	fx.Invoke(InitWorkingDirs),
	fx.Invoke(ValidateEnv),
	fx.Invoke(InitLoggers),
)

func InitWorkingDirs(env *env.Env) error {
	dirs := []string{
		env.GetLogsDir(),
		env.GetUploadsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create dir %s", dir)
		}
	}

	return nil
}

func ValidateEnv(env *env.Env) error {
	return env.Validate()
}

func InitLoggers(
	auditLogger *logger.AuditLogger,
	errorLogger *logger.ErrorLogger,
) error {
	if err := auditLogger.Init(); err != nil {
		return err
	}
	if err := errorLogger.Init(); err != nil {
		return err
	}
	return nil
}
