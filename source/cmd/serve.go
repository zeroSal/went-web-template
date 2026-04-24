package cmd

import (
	"fmt"
	"webtemplate/app"
	"webtemplate/app/config"

	"github.com/zeroSal/went-clio/clio"
	"github.com/zeroSal/went-command/command"

	"github.com/kataras/iris/v12"
)

var _ command.Interface = (*ServeCmd)(nil)
type ServeCmd struct {
	command.Base
}

func NewServeCmd() command.Interface {
	return &ServeCmd{}
}

func (c *ServeCmd) GetHeader() command.Header {
	return command.Header{
		Use:   "serve",
		Short: "Run the web server",
		Long:  "Run the web server on the configured host and port.",
	}
}

func (c *ServeCmd) Invoke() any {
	return c.run
}

func (c *ServeCmd) run(
	env *config.Env,
	clio *clio.Clio,
	buildSpec *app.BuildSpecs,
	irisApp *iris.Application,
) error {
	clio.Banner(buildSpec.GetVersion(), buildSpec.GetChannel(), buildSpec.GetBuildDate())

	addr := fmt.Sprintf("%s:%d", env.Host, env.Port)
	clio.Info("The web server is running on http://" + addr)
	if err := irisApp.Listen(addr, iris.WithoutStartupLog); err != nil {
		return err
	}

	return nil
}
