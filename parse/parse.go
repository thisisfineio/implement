package parse

import (
	"go/token"
	"go/parser"
	"go/ast"
	"io/ioutil"
	"github.com/thisisfineio/implement"
	"strings"
	//"reflect"
)

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
func Inspect(f *ast.File, data []byte) (map[string][]*implement.FunctionSignature, error) {
	signatures := make(map[string][]*implement.FunctionSignature)
	var lastIdent string
	ast.Inspect(f, func(n ast.Node) bool {
		//fmt.Println(n)
		switch  t := n.(type) {
		case *ast.Ident:
			lastIdent = t.Name
		// top level interface definition
		case *ast.InterfaceType:
			sigs := make([]*implement.FunctionSignature, 0)
			for _, f := range t.Methods.List {
				name := getFunctionName(f.Pos(), f.End(), data)
				sig := FunctionSignature(f.Type)

				sig.Name = name
				sigs = append(sigs, sig)
			}
			signatures[lastIdent] = sigs
		}
		return true
	})
	return signatures, nil
}

func getFunctionName(start, end token.Pos, data []byte) string {
	return strings.Split(string(data[start - 1: end -1]), "(")[0]
}

func FunctionSignature(expr ast.Expr) (*implement.FunctionSignature) {
	signature := &implement.FunctionSignature{}
	switch n := expr.(type) {
	// the top level function
	case *ast.FuncType:
		if n.Params != nil {
			//letterMap := make(map[string]int)
			for _, p := range n.Params.List {
				s := getTypeIdentifier(p.Type)
				for _, n := range p.Names {
					param := &implement.Parameter{Name: n.Name, Type: s}
					signature.Parameters = append(signature.Parameters, param)
				}

				if len(p.Names) == 0 {
					param := &implement.Parameter{Type: s}
					signature.Parameters = append(signature.Parameters, param)
				}

			}
		}

		if n.Results != nil {
			for _, r := range n.Results.List {
				s := getTypeIdentifier(r.Type)
				for _, n := range r.Names {
					result := &implement.ReturnValue{Name: n.Name, Type: s}
					signature.ReturnValues = append(signature.ReturnValues, result)
				}

				if len(r.Names) == 0 {
					result := &implement.ReturnValue{Type: s}
					signature.ReturnValues = append(signature.ReturnValues, result)
				}
			}
		}
	}
	return signature
}

func lowerFirstLetterOfVar(s string) string {
	split := strings.Split(s, "*")
	var firstLetter string
	if len(split) > 1 {
		firstLetter = strings.ToLower(split[1])
	} else {
		firstLetter = strings.ToLower(split[0])
	}
	return firstLetter
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
	}
	return ""
}