package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

type Declaration struct {
	Label    string        `json:"label"`
	Type     string        `json:"type"`
	Start    token.Pos     `json:"start"`
	End      token.Pos     `json:"end"`
	Children []Declaration `json:"children,omitempty"`
}

var (
	file = flag.String("f", "", "the path to the file to outline")
)

func main() {
	flag.Parse()
	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, *file, nil, parser.ParseComments)
	if err != nil {
		reportError(fmt.Errorf("Could not parse file %s", *file))
	}

	declarations := []Declaration{}

	for _, decl := range fileAst.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			declarations = append(declarations, Declaration{
				decl.Name.String(),
				"function",
				decl.Pos(),
				decl.End(),
				[]Declaration{},
			})
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.ImportSpec:
					declarations = append(declarations, Declaration{
						spec.Path.Value,
						"import",
						spec.Pos(),
						spec.End(),
						[]Declaration{},
					})
				case *ast.TypeSpec:
					//TODO: Members if it's a struct or interface type?
					declarations = append(declarations, Declaration{
						spec.Name.String(),
						"type",
						spec.Pos(),
						spec.End(),
						[]Declaration{},
					})
				case *ast.ValueSpec:
					for _, id := range spec.Names {
						declarations = append(declarations, Declaration{
							id.Name,
							"variable",
							id.Pos(),
							id.End(),
							[]Declaration{},
						})
					}
				default:
					reportError(fmt.Errorf("Unknown token type: %s", decl.Tok))
				}
			}
		default:
			reportError(fmt.Errorf("Unknown declaration @", decl.Pos()))
		}
	}

	pkg := []*Declaration{&Declaration{
		fileAst.Name.String(),
		"package",
		fileAst.Pos(),
		fileAst.End(),
		declarations,
	}}

	str, _ := json.Marshal(pkg)
	fmt.Println(string(str))

}

func reportError(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
}
