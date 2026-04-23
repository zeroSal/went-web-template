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
			kernel.Run(instance.Invoke())
		},
	).Execute(); err != nil {
		panic(err)
	}
}
