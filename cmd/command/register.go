package command

import (
	"reflect"

	"github.com/spf13/cobra"
)

func Register(
	constructors []func() Interface,
	rootName string,
	run func(instance Interface),
) *cobra.Command {
	root := &cobra.Command{
		Use:   rootName,
		Short: rootName,
		Long:  rootName,
	}

	for _, constructor := range constructors {
		command := constructor()

		cobraCmd := &cobra.Command{
			Use:   command.GetHeader().Use,
			Short: command.GetHeader().Short,
			Long:  command.GetHeader().Long,
		}

		instance := constructor()
		instanceValue := reflect.ValueOf(instance).Elem()

		if header := instance.GetHeader(); header.Flags != nil {
			if flags := header.Flags; flags != nil {
				for _, boolFlag := range flags.Bool {
					fieldName := toFieldName(boolFlag.Name)
					if field, ok := findField(instanceValue, fieldName); ok {
						cobraCmd.Flags().BoolVar(
							field.Interface().(*bool),
							boolFlag.Name,
							boolFlag.Default,
							boolFlag.Usage,
						)
					}
				}
				for _, intFlag := range flags.Int {
					fieldName := toFieldName(intFlag.Name)
					if field, ok := findField(instanceValue, fieldName); ok {
						cobraCmd.Flags().IntVar(
							field.Interface().(*int),
							intFlag.Name,
							intFlag.Default,
							intFlag.Usage,
						)
					}
				}
				for _, stringFlag := range flags.String {
					fieldName := toFieldName(stringFlag.Name)
					if field, ok := findField(instanceValue, fieldName); ok {
						cobraCmd.Flags().StringVar(
							field.Interface().(*string),
							stringFlag.Name,
							stringFlag.Default,
							stringFlag.Usage,
						)
					}
				}
			}
		}

		arguments := instance.GetHeader().Arguments
		cobraCmd.Args = func(cmd *cobra.Command, args []string) error {
			for i, arg := range arguments {
				if i >= len(args) {
					continue
				}
				fieldName := toFieldName(arg.Name)
				if field, ok := findField(instanceValue, fieldName); ok {
					if field.Kind() == reflect.Pointer && field.Type().Elem() == reflect.TypeFor[string]() {
						ptr := args[i]
						field.Set(reflect.ValueOf(&ptr))
					}
				}
			}
			return nil
		}

		cobraCmd.Run = func(cmd *cobra.Command, args []string) {
			run(instance)
		}

		root.AddCommand(cobraCmd)
	}

	return root
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
		if field.Kind() == reflect.Pointer && field.Elem().Kind() == reflect.Struct {
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
