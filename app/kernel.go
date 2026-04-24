package app

import (
	"context"
	"embed"
	"template/app/bootstrap"

	"github.com/zeroSal/went-clio/clio"
	"go.uber.org/fx"
)

type Kernel struct {
	EmbedFS    embed.FS
	BuildSpecs *BuildSpecs
	Clio       *clio.Clio
}

func NewKernel(
	embedFS embed.FS,
	buildSpecs *BuildSpecs,
	clio *clio.Clio,
) *Kernel {
	return &Kernel{
		embedFS,
		buildSpecs,
		clio,
	}
}

func (a *Kernel) Run(invoke any, opts ...fx.Option) error {
	di := []fx.Option{
		Container,
		bootstrap.Init,
		fx.Supply(a.Clio),
		fx.Supply(a.EmbedFS),
		fx.Supply(a.BuildSpecs),
		fx.Invoke(invoke),
		fx.NopLogger,
	}

	app := fx.New(append(di, opts...)...)

	if err := app.Start(context.Background()); err != nil {
		return err
	}

	if err := app.Stop(context.Background()); err != nil {
		return err
	}

	return nil
}
