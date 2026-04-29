package controller

import (
	"webtemplate/app/registry"

	"github.com/kataras/iris/v12"
	"github.com/zeroSal/went-web/controller"
)

var _ controller.Interface = (*Login)(nil)

type Login struct {
	controller.Base
}

func init() {
	registry.Controller.Register(&Login{})
}

func (c *Login) GetConfiguration() controller.Configuration {
	return controller.Configuration{
		Name:  "Login",
	}
}

func (c *Login) Index(ctx iris.Context) {
	ctx.View("login.html.django")
}

func (c *Login) Login(ctx iris.Context) {
	username := ctx.PostValue("username")
	password := ctx.PostValue("password")

	if username == "test" && password == "test" {
		ctx.SetCookieKV("SESSION_ID", username)
		ctx.Redirect("/", iris.StatusFound)
		return
	}

	ctx.ViewData("error", "Invalid credentials")
	ctx.View("login.html.django")
}

func (c *Home) Logout(ctx iris.Context) {
	ctx.RemoveCookie("SESSION_ID")
	ctx.Redirect("/login")
}