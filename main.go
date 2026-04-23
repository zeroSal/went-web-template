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

var EmbedFS embed.FS

func main() {
	buildSpecs := app.NewBuildSpecs(Version, Channel, BuildDate)
	kernel := app.NewKernel(EmbedFS, buildSpecs)

	commands := []command.RegistrableCommand{
		{Name: "serve", Cmd: func() command.Interface { return cmd.NewServeCmd(kernel) }},
		{Name: "greet", Cmd: func() command.Interface { return cmd.NewGreetCmd(kernel) }},
	}

	if err := command.Register(
		commands,
		"template",
		func(instance command.Interface, args []string, flags map[string]any) {
			kernel.Run(instance.Invoke(), args)
		},
	).Execute(); err != nil {
		panic(err)
	}
}