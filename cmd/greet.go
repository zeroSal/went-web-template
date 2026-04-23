package cmd

import (
	"fmt"
	"template/app"
	"template/cmd/command"

	"github.com/zeroSal/go-semantic-log/logger"
)

var _ command.Interface = (*GreetCmd)(nil)

type GreetCmd struct {
	NoWarning *bool
	Name      *string
	Args      *string
}

func NewGreetCmd() command.Interface {
	return &GreetCmd{
		NoWarning: new(bool),
		Name:      new(string),
		Args:      new(string),
	}
}

func (c *GreetCmd) GetHeader() command.Header {
	return command.Header{
		Use:   "greet",
		Short: "Greets the user by name",
		Long:  "Greets the user by name",
		Flags: &command.Flags{
			Bool:   []command.BoolFlag{{Name: "no-warning", Default: false, Usage: "Suppress warning message"}},
			String: []command.StringFlag{{Name: "name", Default: "World", Usage: "Name to greet"}},
		},
	}
}

func (c *GreetCmd) Invoke() any {
	return c.run
}

func (c *GreetCmd) run(
	logger *logger.ConsoleLogger,
	buildSpec *app.BuildSpecs,
) error {
	name := ""
	if c.Args != nil && *c.Args != "" {
		name = *c.Args
	} else if c.Name != nil && *c.Name != "" {
		name = *c.Name
	} else {
		name = "World"
	}
	fmt.Printf("Ciao %s\n", name)

	if buildSpec != nil {
		fmt.Printf("BuildSpec: %v\n", buildSpec.GetBuildDate())
	}

	if c.NoWarning == nil || !*c.NoWarning {
		logger.Warn("WORKSS!")
	}

	return nil
}
