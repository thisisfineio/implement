package implement

import (
	"fmt"
	"bytes"
	"strings"
	"go/format"
	"io/ioutil"
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


const (
	PointerType = iota
	FuncType
	ValueType
)

type Type struct {
	T int
	Value string
}

type Interface struct {
	Name string
	ImplementedName string
	Functions []*FunctionSignature
	NameOptions *NameOptions
}

func (i *Interface) Implement() string {
	buf := bytes.NewBuffer([]byte{})

	// write struct name
	buf.WriteString(fmt.Sprintf("type %s struct {}\n\n", i.ImplementedName))
	for _, f := range i.Functions {
		if len(i.ImplementedName) == 0 {
			i.ImplementedName = " "
		}
		buf.WriteString(fmt.Sprintf("func(%s *%s) %s(", strings.ToLower(string(i.ImplementedName[0])), i.ImplementedName, f.Name))
		for index, p := range f.Parameters {
			if index == len(f.Parameters) - 1 {
				buf.WriteString(p.NameT(i.NameOptions) + " " + p.Type)
			} else {
				buf.WriteString(p.NameT(i.NameOptions) + " " + p.Type + ", ")
			}

		}
		buf.WriteString(") ")
		if len(f.ReturnValues) > 1 {
			buf.WriteString("(")
		}
		for i, r := range f.ReturnValues {

			if i == len(f.ReturnValues) - 1 {
				buf.WriteString(r.Name + " " + r.Type)
			} else {
				buf.WriteString(r.Name + " " + r.Type + ", ")
			}
		}

		if len(f.ReturnValues) > 1 {
			buf.WriteString(")")
		}
		buf.WriteString("{\n")

		for i, r := range f.ReturnValues {
			if i == 0 {
				buf.WriteString("return ")
			}

			if i == len(f.ReturnValues) - 1 {
				buf.WriteString(ZeroValueString(r.Type))
			} else {
				buf.WriteString(ZeroValueString(r.Type) + ", ")
			}

		}
		buf.WriteString("\n}\n")
	}
	return buf.String()
}

func (i *Interface) Save(filepath string) error {
	data := i.Data()
	ioutil.WriteFile(filepath, data, )
}

func ZeroValueString(s string) string {
	switch {
	case strings.Contains(s, "func"):
		return "nil"
	case strings.Contains(s, "*"):
		return "nil"
	case strings.Contains(s, "int"),
		strings.Contains(s, "float"),
		strings.Contains(s, "complex"),
		strings.Contains(s, "byte"),
		strings.Contains(s, "rune"):
		return "0"
	case strings.Contains(s, "bool"):
		return "false"
	case strings.Contains(s, "string"):
		return "\"\""
	case strings.Contains(s, "error"):
		return "nil"
	default:
		return s + "{}"

	}
}

func (i *Interface) String() string {
	data, _ := format.Source([]byte(i.Implement()))
	return string(data)
}

func (i *Interface) Data() []byte {
	return []byte(i.String())
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

func LowerFirstLetterOfVar(s string) string {
	split := strings.Split(s, ".")
	if len(split) > 1 {
		return strings.ToLower(string(split[1][0]))
	}
	split = strings.Split(s, "*")
	if len(split) > 1 {
		return strings.ToLower(string(split[1][0]))
	}
	return strings.ToLower(string(s[0]))
}

func (p *Parameter) NameT(opts *NameOptions) string {
	if p.Name != "" {
		return p.Name
	}
	if opts == nil {
		opts = DefaultNameOptions()
	}
	switch opts.ParameterOpt {

	case FirstLetterAndNumbers:
		return LowerFirstLetterOfVar(p.Type)
	case FullIdentityName:
		return ""
	}
	return ""
}

type ReturnValue struct {
	Type string
	Name string
}