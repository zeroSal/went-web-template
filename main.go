package main

import (
	"embed"
	"template/app"
	"template/cmd"
	"template/cmd/command"
)

var Version = ""
var Channel = ""
var BuildDate = ""

//go:embed templates/* static/*
var EmbedFS embed.FS

func main() {
	buildSpecs := app.NewBuildSpecs(Version, Channel, BuildDate)
	kernel := app.NewKernel(EmbedFS, buildSpecs)

	commands := []func() command.Interface{
		cmd.NewGreetCmd,
		cmd.NewServeCmd,
	}

	if err := command.Register(
		commands,
		"template",
		func(instance command.Interface) {
			if err := kernel.Run(instance.Invoke()); err != nil {
				panic(err)
			}
		},
	).Execute(); err != nil {
		panic(err)
	}
}
