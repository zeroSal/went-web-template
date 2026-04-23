package command

import (
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type RegistrableCommand struct {
	Name   string
	Cmd    func() Interface
}

func Register(
	commands []RegistrableCommand,
	rootName string,
	run func(instance Interface, args []string, flagValues map[string]any),
) *cobra.Command {
	root := &cobra.Command{
		Use:   rootName,
		Short: rootName,
		Long:  rootName,
	}

	for _, c := range commands {
		constructor := c.Cmd
		name := c.Name

		cobraCmd := &cobra.Command{
			Use:   name,
			Short: "Command " + name,
			Long:  "Command " + name,
		}

		instance := constructor()
		instanceValue := reflect.ValueOf(instance).Elem()

		if header := instance.GetHeader(); header.Flags != nil {
			if flags := header.Flags; flags != nil {
				for _, boolFlag := range flags.Bool {
					fieldName := toFieldName(boolFlag.Name)
					if field, ok := findField(instanceValue, fieldName); ok {
						cobraCmd.Flags().BoolVar(field.Interface().(*bool), boolFlag.Name, boolFlag.Default, boolFlag.Usage)
					}
				}
				for _, intFlag := range flags.Int {
					fieldName := toFieldName(intFlag.Name)
					if field, ok := findField(instanceValue, fieldName); ok {
						cobraCmd.Flags().IntVar(field.Interface().(*int), intFlag.Name, intFlag.Default, intFlag.Usage)
					}
				}
				for _, stringFlag := range flags.String {
					fieldName := toFieldName(stringFlag.Name)
					if field, ok := findField(instanceValue, fieldName); ok {
						cobraCmd.Flags().StringVar(field.Interface().(*string), stringFlag.Name, stringFlag.Default, stringFlag.Usage)
					}
				}
			}
		}

		cobraCmd.Run = func(cmd *cobra.Command, args []string) {
			parsedFlags := parseFlags(cmd)
			run(instance, args, parsedFlags)
		}

		root.AddCommand(cobraCmd)
	}

	return root
}

func parseFlags(cmd *cobra.Command) map[string]any {
	flags := make(map[string]any)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		switch f.Value.Type() {
		case "bool":
			flags[f.Name], _ = cmd.Flags().GetBool(f.Name)
		case "int":
			flags[f.Name], _ = cmd.Flags().GetInt(f.Name)
		case "string":
			flags[f.Name], _ = cmd.Flags().GetString(f.Name)
		}
	})
	return flags
}

func toFieldName(name string) string {
	result := make([]byte, 0, len(name))
	for i, c := range name {
		if c == '-' {
			continue
		}
		if i == 0 || name[i-1] == '-' {
			result = append(result, byte(c-'a'+'A'))
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}

func findField(v reflect.Value, name string) (reflect.Value, bool) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Ptr && field.Elem().Kind() == reflect.Struct {
			if found, ok := findField(field.Elem(), name); ok {
				return found, true
			}
		}
		if v.Type().Field(i).Name == name {
			return field, true
		}
	}
	return reflect.Value{}, false
}