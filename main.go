package main

import (
	"fmt"
	"os"
	"go/ast"
	"go/token"
	"go/parser"
	"flag"
	"encoding/json"
)

type Declaration struct {
	Label		string 			`json:"label"`
	Type		string 			`json:"type"`
	Icon		string 			`json:"icon,emitempty"`
	Start		token.Pos    	`json:"start"`
	End   		token.Pos    	`json:"end"`
	Children 	[]Declaration	`json:"children"`
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
				"",
				"function",
				decl.Pos(),
				decl.End(),
				[]Declaration{},
			})
		case *ast.GenDecl:
			switch decl.Tok {
			case token.IMPORT:
				for _, spec := range(decl.Specs) {
					importSpec := spec.(*ast.ImportSpec)
					declarations = append(declarations, Declaration{
						importSpec.Path.Value,
						"",
						"import",
						importSpec.Pos(),
						importSpec.End(),
						[]Declaration{},
					})
				}
			case token.TYPE:
				fmt.Printf("TYPE: %s [%s]\n", decl.Specs)
			case token.VAR:
				fmt.Printf("VAR: %s [%s]\n",  decl.Specs)
			case token.CONST:
				fmt.Printf("CONST: %s [%s]\n", decl.Specs)
			default:
				reportError(fmt.Errorf("Unknown token type: %s", decl.Tok))
			}
		}
	}
	
	pkg := &Declaration{
		fileAst.Name.String(),
		"",
		"package",
		fileAst.Pos(),
		fileAst.End(),
		declarations,
	}
	
	str, _ := json.MarshalIndent(pkg,"", " ")
	fmt.Println(string(str))
	
}

func reportError(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
}
