package registry

import (
	"github.com/zeroSal/went-command/command"
	"github.com/zeroSal/went-web/controller"
)

var Command = command.NewRegistry()
var Controller = controller.NewRegistry()