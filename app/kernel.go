package app

import (
	"context"
	"embed"
	"errors"
	"os"
	"template/app/config"
	"template/di"

	"github.com/kataras/iris/v12"
	"go.uber.org/fx"
)

type Kernel struct {
	EmbedFS    embed.FS
	BuildSpecs *BuildSpecs
}

func NewKernel(
	embedFS embed.FS,
	buildSpecs *BuildSpecs,
) *Kernel {
	return &Kernel{
		embedFS,
		buildSpecs,
	}
}

func (a *Kernel) Run(invoke any, opts ...fx.Option) error {
	buildSpec := func() *BuildSpecs {
		return a.BuildSpecs
	}

	appOpts := []fx.Option{
		di.Container,
		fx.Supply(a.EmbedFS),
		fx.Provide(buildSpec),
		fx.Provide(config.LoadEnv),
		fx.Provide(InitIris),
		fx.Invoke(InitWorkingDirs),
		fx.Invoke(invoke),
		fx.NopLogger,
	}

	app := fx.New(append(appOpts, opts...)...)

	if err := app.Start(context.Background()); err != nil {
		return err
	}

	if err := app.Stop(context.Background()); err != nil {
		return err
	}

	return nil
}

func InitWorkingDirs(
	env *config.Env,
) error {
	dirs := []string{
		env.GetLogsDir(),
		env.GetUploadsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.New("Cannot create dir " + dir)
		}
	}

	return nil
}

func InitIris(env *config.Env, embedFS embed.FS) *iris.Application {
	app := iris.New()

	engine := iris.Django(embedFS, ".html.django")
	if env.Env == "dev" {
		engine.Reload(true)
	}

	app.RegisterView(engine)
	app.HandleDir("/static", embedFS)

	return app
}
