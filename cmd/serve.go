package cmd

import (
	"fmt"
	"template/app"
	"template/app/config"
	"template/cmd/command"

	"github.com/kataras/iris/v12"
	"github.com/zeroSal/go-semantic-log/logger"
)

var _ command.Interface = (*GreetCmd)(nil)

type ServeCmd struct {
}

func NewServeCmd() command.Interface {
	return &ServeCmd{}
}

func (c *ServeCmd) GetHeader() command.Header {
	return command.Header{
		Use:   "serve",
		Short: "Run the web server",
		Long:  "Run the web server",
	}
}

func (c *ServeCmd) Invoke() any {
	return c.run
}

func (c *ServeCmd) run(
	env *config.Env,
	log *logger.ConsoleLogger,
	buildSpec *app.BuildSpecs,
	irisApp *iris.Application,
) error {
	addr := fmt.Sprintf("%s:%d", env.Host, env.Port)
	log.Info("The web server is running on http://" + addr)
	log.Info(buildSpec.GetBuildDate())
	if err := irisApp.Listen(addr, iris.WithoutStartupLog); err != nil {
		return err
	}

	return nil
}
