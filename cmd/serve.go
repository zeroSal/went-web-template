package cmd

import (
	"fmt"
	"template/app"
	"template/app/config"
	"template/cmd/command"

	"github.com/kataras/iris/v12"
	"github.com/zeroSal/go-semantic-log/logger"
)

type ServeCmd struct {
	kernel *app.Kernel
}

func NewServeCmd(kernel *app.Kernel) command.Interface {
	return &ServeCmd{kernel: kernel}
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
	irisApp *iris.Application,
) {
	addr := fmt.Sprintf("%s:%d", env.Host, env.Port)
	log.Info("The web server is running on http://" + addr)
	irisApp.Listen(addr, iris.WithoutStartupLog)
}