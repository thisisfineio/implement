package implement


type Options struct {
	InputPath string
	OutputPath string
	OverwriteExisting bool
	Interfaces []string
	Implementors []string
}


func (o *Options) Validate() error {
	if len(o.Interfaces) != len(o.Implementors) {
		return ErrInvalidLengths
	}
	return nil
}

const (
	Func = "func"
	Struct = "struct"
	Interface = "interface"
)

type FunctionSignature struct {
	Name string
	Parameters []*Parameter
	ReturnValues []*ReturnValue
}

func (f *FunctionSignature) String() string {
return ""
}

type Parameter struct {
	Type string
	Name string
}

type ReturnValue struct {
	Type string
	Name string
}