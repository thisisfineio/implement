package parse

import (
	"go/token"
	"go/parser"
	"go/ast"
	"fmt"
	"io/ioutil"
	"github.com/thisisfineio/implement"
	"strings"
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

	ast.Inspect(f, func(n ast.Node) bool {
		//fmt.Println(n)
		switch  t := n.(type) {
		// top level interface definition
		case *ast.InterfaceType:
			for _, f := range t.Methods.List {
				//fmt.Println(f)
				FunctionSignature(f.Type)
			}
		}
		return true
	})
	return []byte{}, nil
}


func FunctionSignature(expr ast.Expr) *implement.FunctionSignature {
	//signature := &implement.FunctionSignature{}
	switch n := expr.(type) {
	// the top level function
	case *ast.FuncType:

		if n.Params != nil {
			fmt.Println("Params")
			//letterMap := make(map[string]int)
			for _, p := range n.Params.List {
				s, err := getTypeIdentifier(p.Type)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println(s)
				var typ string

				for _, n := range p.Names {
					param := &implement.Parameter{Name: n.Name}
					param.Type = typ
					//fmt.Println(param)
				}

				if len(p.Names) == 0 {
					//name := getVarName()
					param := &implement.Parameter{Type: typ}
					param.Type = typ
				}

			}
			fmt.Println()
		}

		if n.Results != nil {
			fmt.Println("Returns")
			for _, r := range n.Results.List {
				//result := &implement.ReturnValue{}
				//fmt.Println(r.Names)
				s, err := getTypeIdentifier(r.Type)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(s)
			}
			fmt.Println()
		}
	}
	return nil
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

func parseInterfaceType(i *ast.InterfaceType){

}

func parseFuncType(f *ast.FuncType) {

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