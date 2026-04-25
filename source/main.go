package main

import (
	"embed"
	"fmt"
	"os"
	"webtemplate/app"
	"webtemplate/cmd"

	"github.com/zeroSal/went-clio/clio"
	"github.com/zeroSal/went-command/command"

	"github.com/spf13/cobra"
)

var Version = ""
var Channel = ""
var BuildDate = ""

//go:embed templates/* static/*
var EmbedFS embed.FS

func main() {
	clio := clio.NewClio()

	data, err := EmbedFS.ReadFile("templates/banner.template")
	if err != nil {
		clio.Error("Error loading the banner template.")
		os.Exit(3)
	}

	buildSpecs := app.NewBuildSpecs(Version, Channel, BuildDate)
	clio.SetBannerTemplate(string(data))

	kernel := app.NewKernel(EmbedFS, buildSpecs, clio)

	commands := []func() command.Interface{
		cmd.NewServeCmd,
	}

	root := &cobra.Command{
		Version: fmt.Sprintf("%s-%s (%s)", Version, Channel, BuildDate),
		Use:   "webtemplate",
		Short: "{{ SHORT_PROJECT_DESCRIPTION }}",
		Long:  "{{ LONG_PROJECT_DESCRIPTION }}",
	}

	run := func(instance command.Interface) {
		if err := kernel.Run(instance.Invoke()); err != nil {
			clio.Fatal(err.Error())
			os.Exit(1)
		}
	}

	if err := command.Register(commands, root, run).Execute(); err != nil {
		clio.Fatal("COMMAND REGISTRATION ERROR: " + err.Error())
		os.Exit(2)
	}
}
