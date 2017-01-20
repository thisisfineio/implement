package implement

import (
	"fmt"
	"bytes"
	"strings"
	"go/format"
	"io/ioutil"
	"os"
	"errors"
	"go/ast"
	"go/token"
	"go/parser"
)

type Options struct {
	InputPath string
	OutputPath string
	OverwriteExisting bool
	ImplementorMap map[string]string
	NameOptions *NameOptions
}

func DefaultOptions() *Options {
	return &Options{NameOptions: DefaultNameOptions(), OverwriteExisting: false}
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
	Options *Options
}

func (i *Interface) Implement() string {
	if i.Options == nil {
		i.Options = DefaultOptions()
	}
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
				buf.WriteString(p.NameT(i.Options.NameOptions) + " " + p.Type)
			} else {
				buf.WriteString(p.NameT(i.Options.NameOptions) + " " + p.Type + ", ")
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
		buf.WriteString("\n}\n\n")
	}
	return buf.String()
}

var (
	FileExistsErr = errors.New("File exists. If you would like to overwrite it, provide the OverwriteExisting option.")
)

func (i *Interface) Save(dirpath string) error {
	filepath := dirpath + string(os.PathSeparator) + i.Name + ".go"
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		if !i.Options.OverwriteExisting {
			return FileExistsErr
		}
	}

	data := i.Data()
	split := strings.Split(dirpath, string(os.PathSeparator))
	d := split[len(split) - 1]
	pack := "package " + d + "\n"
	data = append([]byte(pack), data...)
	return ioutil.WriteFile(filepath, data, 0664)

}

func ZeroValueString(s string) string {
	switch {
	case strings.Contains(s, "func"):
		return "nil"
	case strings.Contains(s, "[]"):
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
	data, err := format.Source([]byte(i.Implement()))
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (i *Interface) Data() []byte {
	return []byte(i.String())
}

func Interfaces(signatures map[string][]*FunctionSignature, m map[string]string) []*Interface {
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
	split = strings.Split(s, "[]")
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



// Parse returns an *ast.File, and an error
func Parse(data []byte, filename string) (*ast.File, []byte, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, data, parser.ParseComments)
	return f, data, err
}

func File(filename string) (*ast.File, []byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, err
	}
	return Parse(data, filename)
}

// Inspect visits all nodes in the *ast.File (recursively)
// data is necessary for us to determine interface names
func Inspect(f *ast.File, data []byte) map[string][]*FunctionSignature {
	signatures := make(map[string][]*FunctionSignature)
	var lastIdent string
	ast.Inspect(f, func(n ast.Node) bool {
		//fmt.Println(n)
		switch  t := n.(type) {
		case *ast.Ident:
			lastIdent = t.Name
		// top level interface definition
		case *ast.InterfaceType:
			sigs := make([]*FunctionSignature, 0)
			for _, f := range t.Methods.List {
				name := getFunctionName(f.Pos(), f.End(), data)
				sig := GetFunctionSignatures(f.Type, data)

				sig.Name = name
				sigs = append(sigs, sig)
			}
			signatures[lastIdent] = sigs
		}
		return true
	})
	return signatures
}

func getFunctionName(start, end token.Pos, data []byte) string {
	return strings.Split(string(data[start - 1: end -1]), "(")[0]
}

// TODO - figure out how to actually return this... maybe send finished signatures to a channel instead of returning the function
func GetFunctionSignatures(expr ast.Expr, data []byte) (*FunctionSignature) {
	signature := &FunctionSignature{}
	switch n := expr.(type) {
	// the top level function
	case *ast.FuncType:
		signature.Name = getFunctionName(expr.Pos(), expr.End(), data)
		if n.Params != nil {
			//letterMap := make(map[string]int)
			for _, p := range n.Params.List {
				s := getTypeIdentifier(p.Type)
				for _, n := range p.Names {
					param := &Parameter{Name: n.Name, Type: s}
					signature.Parameters = append(signature.Parameters, param)
				}

				if len(p.Names) == 0 {
					param := &Parameter{Type: s}
					signature.Parameters = append(signature.Parameters, param)
				}

			}
		}

		if n.Results != nil {
			for _, r := range n.Results.List {
				s := getTypeIdentifier(r.Type)
				for _, n := range r.Names {
					result := &ReturnValue{Name: n.Name, Type: s}
					signature.ReturnValues = append(signature.ReturnValues, result)
				}

				if len(r.Names) == 0 {
					result := &ReturnValue{Type: s}
					signature.ReturnValues = append(signature.ReturnValues, result)
				}
			}
		}

	case *ast.Ident:
		if t, ok := n.Obj.Decl.(*ast.TypeSpec); ok {
			switch t := t.Type.(type) {
			case *ast.InterfaceType:
				for _, f := range t.Methods.List {
					sig := GetFunctionSignatures(f.Type, data)
					fmt.Printf("%+v\n",sig)
				}
			}
		}


	default:

	}
	return signature
}

func getTypeIdentifier(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.String()
	case *ast.StarExpr:
		s := getTypeIdentifier(t.X)
		return "*" + s
	case *ast.SelectorExpr:
		s := getTypeIdentifier(t.X)
		return s + "." + t.Sel.String()
	case *ast.FuncType:
		var typ = "func("
		if t.Params != nil {
			for _, p := range t.Params.List {
				s := getTypeIdentifier(p.Type)
				if len(p.Names) > 0 {
					for i := range p.Names {
						typ += s
						if i != len(p.Names) - 1 {
							typ += ","
						}
					}

				} else {
					typ += s
				}
			}
		}
		typ += ")"
		if t.Results != nil {
			if len(t.Results.List) > 1 {
				typ += " ("
			}

			for i, r := range t.Results.List {
				s := getTypeIdentifier(r.Type)
				typ += s
				if len(t.Results.List) > 1 && i != len(t.Results.List) - 1 {
					typ += ", "
				}
			}

			if len(t.Results.List) > 1 {
				typ += ")"
			}
		}
		return typ

	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ArrayType:
		s := getTypeIdentifier(t.Elt)
		return "[]" + s
	default:
	}
	return ""
}