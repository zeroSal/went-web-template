package command

type Header struct {
	Use       string
	Short     string
	Long      string
	Flags     *Flags
	Arguments []Argument
}

type Flags struct {
	Bool   []BoolFlag
	Int    []IntFlag
	String []StringFlag
}

type BoolFlag struct {
	Name    string
	Default bool
	Usage   string
}

type IntFlag struct {
	Name    string
	Default int
	Usage   string
}

type StringFlag struct {
	Name    string
	Default string
	Usage   string
}

type Argument struct {
	Name    string
	Default string
	Usage   string
}
