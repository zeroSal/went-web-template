package ctrl

import (
	"webtemplate/app/service/env"
	"webtemplate/registry"

	"github.com/kataras/iris/v12"
	"github.com/zeroSal/went-web/controller"
)

var _ controller.Interface = (*Home)(nil)
type Home struct {
	controller.Base
	env *env.Env
}

func init() {
	registry.Controller.Register(NewHome)
}

func NewHome(
	env *env.Env,
) *Home {
	return &Home{
		env: env,
	}
}

func (c *Home) Register(app *iris.Application) {
	app.Get("/", c.index())
}

func (c *Home) index() iris.Handler {
	return func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "Hello from Iris!",
			"env":     c.env.Env,
		});
	}
}
