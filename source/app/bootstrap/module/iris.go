package module

import (
	"embed"
	"webtemplate/app/config"

	"github.com/kataras/iris/v12"
)

func IrisProvider(
	env *config.Env,
	embedFS embed.FS,
) *iris.Application {
	app := iris.New()

	engine := iris.Django(embedFS, ".html.django")
	if env.Env == "dev" {
		engine.Reload(true)
	}

	app.RegisterView(engine)
	app.HandleDir("/static", embedFS)

	return app
}