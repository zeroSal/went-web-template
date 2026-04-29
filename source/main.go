package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"webtemplate/app"
	_ "webtemplate/app/cmd"
	_ "webtemplate/app/controller"
	"webtemplate/app/registry"

	"github.com/zeroSal/went-clio/clio"
	"github.com/zeroSal/went-command/command"

	"github.com/spf13/cobra"
)

var Version = ""
var Channel = ""
var BuildDate = ""

//go:embed res/* static/* templates/* config/*
var EmbedFS embed.FS

func main() {
	clio := clio.NewClio()

	data, err := EmbedFS.ReadFile("res/banner.template")
	if err != nil {
		clio.Error("Error loading the banner template.")
		os.Exit(3)
	}

	specs := app.NewSpecs(Version, Channel, BuildDate)
	clio.SetBanner(string(data), Version, Channel, BuildDate)

	ctx := context.Background()

	kernel := app.NewKernel(ctx, EmbedFS, specs, clio)

	root := &cobra.Command{
		Version: fmt.Sprintf("%s-%s (%s)", Version, Channel, BuildDate),
		Use:     "clitemplate",
		Short:   "{{ SHORT_PROJECT_DESCRIPTION }}",
		Long:    "{{ LONG_PROJECT_DESCRIPTION }}",
	}

	run := func(command command.Interface) {
		if err := kernel.Run(command.Invoke()); err != nil {
			clio.Fatal("%s", err.Error())
			os.Exit(1)
		}
	}

	if err := command.Mount(registry.Command.All(), root, run).Execute(); err != nil {
		clio.Fatal("Error mounting commands: %s", err.Error())
		os.Exit(2)
	}
}
