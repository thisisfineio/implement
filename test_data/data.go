package test_data

import (
	"io"
	"bytes"
)

type Struct struct {

}

type TestInterface interface {
	Func()
	ParamFunc(*Struct)
	ReturnFunc() string
	MultiIntParam(i, j int)
	TupleFunc() (string, error)
	MultiParam(s *Struct, reader io.Reader)
	Combo(*Struct, io.Reader) (string, error)
	FunctionParam(f func(i int) (interface{}, error))
	FunctionParamAndReturn(f func(i int) (interface{}, error)) (*Struct, error)
	AnonymousFuncs(func(func(i, j int) func()) func(i int)) func(int) // efffffff recursion hell
	PackageStruct(*bytes.Buffer) *bytes.Reader
}

type IgnoredInterface interface {
	IgnoredFunction(i int) *Struct
}