package cmd

import (
	"context"
	"fmt"
	"webtemplate/app/service/env"
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
		Short: "Start the web server",
		Long:  "Start the web server.",
	}
}

func (c *Serve) Invoke() any {
	return c.run
}

func (c *Serve) run(
	ctx context.Context,
	env *env.Env,
	clio *clio.Clio,
	app *iris.Application,
) error {
	clio.Banner()

	bind := fmt.Sprintf("%s:%d", env.Host, env.Port)

	clio.Info("Server bound on %s", bind)

	return app.Listen(bind, iris.WithoutStartupLog)
}
