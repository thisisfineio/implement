package parse

import (
	"go/token"
	"go/parser"
	"go/ast"
	"fmt"
	"io/ioutil"
	"github.com/thisisfineio/implement"
	"strings"
	"unicode"
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
func Inspect(f *ast.File, data []byte) ([]byte, error) {
	fmt.Println()
	signatures := make(map[string][]*implement.FunctionSignature)
	ast.Inspect(f, func(n ast.Node) bool {
		//fmt.Println(n)
		switch  t := n.(type) {
		// top level interface definition
		case *ast.InterfaceType:
			interfaceName := getInterfaceName(t.Pos(), data)
			sigs := make([]*implement.FunctionSignature, 0)
			for i, f := range t.Methods.List {
				fmt.Println("Function", i)
				//fmt.Println(f)
				name := getFunctionName(f.Pos(), f.End(), data)
				fmt.Println(name)
				sig := FunctionSignature(f.Type)
				sig.Name = name
				sigs = append(sigs, sig)
			}
			signatures[interfaceName] = sigs
		}
		return true
	})
	for k, v := range signatures {
		fmt.Println("Interface:", k)
		for _, f := range v {
			fmt.Println(f)
		}
	}
	return []byte{}, nil
}

func getFunctionName(start, end token.Pos, data []byte) string {
	return strings.Split(string(data[start - 1: end -1]), "(")[0]
}

// todo - does not work, need to scan backwards from start to find the second space
func getInterfaceName(start token.Pos, data []byte) string {
	cur := start
	spaceFound := false
	for cur != -1 {
		if unicode.IsSpace(rune(data[cur])) {
			if spaceFound {
				return string(data[cur+1:start -1])
				break
			}
			spaceFound = true
		}
		cur--
	}
	return ""
}

func FunctionSignature(expr ast.Expr) *implement.FunctionSignature {
	signature := &implement.FunctionSignature{}
	switch n := expr.(type) {
	// the top level function
	case *ast.FuncType:
		if n.Params != nil {
			//letterMap := make(map[string]int)
			for _, p := range n.Params.List {
				s, err := getTypeIdentifier(p.Type)
				if err != nil {
					fmt.Println(err)
				}
				for _, n := range p.Names {
					param := &implement.Parameter{Name: n.Name, Type: s}
					signature.Parameters = append(signature.Parameters, param)
				}

				if len(p.Names) == 0 {
					param := &implement.Parameter{Type: s}
					signature.Parameters = append(signature.Parameters, param)
				}

			}
			fmt.Println()
		}

		if n.Results != nil {
			for _, r := range n.Results.List {
				s, err := getTypeIdentifier(r.Type)
				if err != nil {
					fmt.Println(err)
				}
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

func getTypeIdentifier(expr ast.Expr) (string, error) {

	switch t := expr.(type) {
	case *ast.Ident:
		return t.String(), nil
	case *ast.StarExpr:
		s, err := getTypeIdentifier(t.X)
		if err != nil {
			return "", err
		}
		return "*" + s,  nil
	case *ast.SelectorExpr:
		s, err := getTypeIdentifier(t.X)
		if err != nil {
			return "", err
		}
		return s + "." + t.Sel.String(), nil
	case *ast.FuncType:
		var typ = "func("
		if t.Params != nil {
			for _, p := range t.Params.List {
				s, err := getTypeIdentifier(p.Type)
				if err != nil {
					return "", err
				}
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
				s, err := getTypeIdentifier(r.Type)
				if err != nil {
					return "", err
				}
				typ += s
				if len(t.Results.List) > 1 && i != len(t.Results.List) - 1 {
					typ += ", "
				}
			}

			if len(t.Results.List) > 1 {
				typ += ")"
			}
		}
		return typ, nil

	case *ast.InterfaceType:
		return "interface{}", nil
	}

	return "", nil
}

//func getVarName(o *implement.Options, letterMap map[string]int) string {
//
//}

//func Params(buffer string) ([]*implement.Parameter, error) {
//	split := strings.Split(buffer, ",")
//
//
//}
//
//func Returns(buffer string) ([]*implement.ReturnValue, error) {
//
//}