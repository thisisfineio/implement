package implement

import (
	"fmt"
	"bytes"
	"strings"
)

type Options struct {
	InputPath string
	OutputPath string
	OverwriteExisting bool
	ImplementorMap map[string]string
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

type Interface struct {
	Name string
	ImplementedName string
	Functions []*FunctionSignature
}

// TODO - finish this!!!
func (i *Interface) Implement() string {
	buf := bytes.NewBuffer([]byte{})

	// write struct name
	buf.WriteString(fmt.Sprintf("type %s struct {}\n\n", i.ImplementedName))
	for _, f := range i.Functions {
		if len(i.ImplementedName) == 0 {
			i.ImplementedName = " "
		}
		buf.WriteString(fmt.Sprintf("func(%s *%s) %s(", strings.ToLower(string(i.ImplementedName[0])), i.ImplementedName, f.Name))
		for _, p := range f.Parameters {
			buf.WriteString(p.Name + " " + p.Type)
		}
		buf.WriteString(") ")

		for _, r := range f.ReturnValues {
			buf.WriteString(r.Name + " " + r.Type)
		}
		buf.WriteString("{\n}\n")
	}
	return buf.String()
}

func (i *Interface) String() string {
	return i.Implement()
}

func GetInterfaces(signatures map[string][]*FunctionSignature, m map[string]string) []*Interface {
	interfaces := make([]*Interface, 0)
	for k, v := range signatures {
		if _, ok := m[k]; !ok {
			continue
		}
		i := &Interface{Name: k, ImplementedName: m[k], Functions: v}
		interfaces = append(interfaces, i)
	}
	return interfaces
}

type FunctionSignature struct {
	Name string
	Parameters []*Parameter
	ReturnValues []*ReturnValue
}

type Parameter struct {
	Type string
	Name string
}

type ReturnValue struct {
	Type string
	Name string
}