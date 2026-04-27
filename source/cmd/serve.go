package cmd

import (
	"context"
	"fmt"
	"webtemplate/app/config"
	"webtemplate/registry"

	"github.com/kataras/iris/v12"
	"github.com/zeroSal/went-clio/clio"
	"github.com/zeroSal/went-command/command"
)

var _ command.Interface = (*Serve)(nil)
type Serve struct {
	command.Base
}

func init() {
    registry.Command.Register(&Serve{})
}

func (c *Serve) GetHeader() command.Header {
	return command.Header{
		Use:   "serve",
		Short: "Starts the web server",
		Long:  "Starts the web server on the configured host and port.",
	}
}

func (c *Serve) Invoke() any {
	return c.run
}

func (c *Serve) run(
	ctx context.Context,
	env *config.Env,
	clio *clio.Clio,
	irisApp *iris.Application,
) error {
	clio.Banner()

	addr := fmt.Sprintf("%s:%d", env.Host, env.Port)
	clio.Info("The web server is listening on http://%s", addr)
	if err := irisApp.Listen(addr, iris.WithoutStartupLog); err != nil {
		return err
	}

	return nil
}
