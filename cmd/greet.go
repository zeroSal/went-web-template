package cmd

import (
	"fmt"
	"template/app"
	"template/cmd/command"

	"github.com/zeroSal/go-semantic-log/logger"
)

type GreetCmd struct {
	kernel     *app.Kernel
	NoWarning  *bool
	Name       *string
}

func NewGreetCmd(kernel *app.Kernel) command.Interface {
	return &GreetCmd{
		kernel:     kernel,
		NoWarning:  new(bool),
		Name:       new(string),
	}
}

func (c *GreetCmd) GetHeader() command.Header {
	return command.Header{
		Use:   "greet",
		Short: "Greets the user by name",
		Long:  "Greets the user by name",
		Flags: &command.Flags{
			Bool: []command.BoolFlag{{Name: "no-warning", Default: false, Usage: "Suppress warning message"}},
		},
		Arguments: []command.Argument{{Name: "name", Default: "World", Usage: "Name to greet"}},
	}
}

func (c *GreetCmd) Invoke() any {
	return c.run
}

func (c *GreetCmd) run(args []string, logger *logger.ConsoleLogger) {
	name := *c.Name
	if len(args) > 0 {
		name = args[0]
	}
	fmt.Printf("Ciao %s\n", name)

	if c.NoWarning == nil || !*c.NoWarning {
		logger.Warn("WORKSS!")
	}
}