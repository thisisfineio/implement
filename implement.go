package implement

import "fmt"

type Options struct {
	InputPath string
	OutputPath string
	OverwriteExisting bool
	Interfaces []string
	Implementors []string
	NameOptions *NameOptions
}

const (
	FirstLetterAndNumbers = iota
	FullIdentityName
)

type NameOptions struct {
	ParameterOpt int
}

func DefaultNameOptions() *NameOptions {
	return &NameOptions{FirstLetterAndNumbers}
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
	NameOptions *NameOptions

}

func (f *FunctionSignature) String() string {

}

type Parameter struct {
	Type string
	Name string
}

type ReturnValue struct {
	Type string
	Name string
}