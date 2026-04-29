package controller

import (
	"webtemplate/app/registry"

	"github.com/kataras/iris/v12"
	"github.com/zeroSal/went-web/controller"
)

var _ controller.Interface = (*Home)(nil)

type Home struct {
	controller.Base
}

func init() {
	registry.Controller.Register(&Home{})
}

func (c *Home) GetConfiguration() controller.Configuration {
	return controller.Configuration{
		Name:  "Home",
	}
}

func (c *Home) Index(ctx iris.Context) {
	ctx.View("home.html.django")
}
