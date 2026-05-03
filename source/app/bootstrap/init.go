package bootstrap

import (
	"fmt"
	"os"
	"webtemplate/app/service/env"
	"webtemplate/app/service/logger"

	"github.com/zeroSal/went-web/controller"

	"go.uber.org/fx"
)

var Init = fx.Options(
	fx.Invoke(initWorkingDirs),
	fx.Invoke(validateEnv),
	fx.Invoke(initLoggers),
	fx.Invoke(controller.Mount),
)

func initLoggers(
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

func initWorkingDirs(env *env.Env) error {
	dirs := []string{
		env.GetLogsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create dir %s", dir)
		}
	}

	return nil
}

func validateEnv(env *env.Env) error {
	return env.Validate()
}
