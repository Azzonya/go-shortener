// Package exit provides analyzer for no os exit calls inside main-like functions
package exit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"strings"
)

// New create new instance of osexitcheck
func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "osexitcheck",
		Doc:  "checks for direct os.Exit calls inside main-like functions",
		Run:  run,
	}
}

// run stars osexitcheck analyzer
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if file.Name.String() != "main" {
			continue
		}
		if !strings.HasSuffix(pass.Fset.Position(file.Pos()).Filename, ".go") {
			continue
		}
		ast.Inspect(file, inspect(pass))
	}
	return nil, nil
}

// inspect search for main function and check if os.Exit() is called.
func inspect(pass *analysis.Pass) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		callExpr, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		if funcLit, ok := callExpr.Fun.(*ast.FuncLit); ok {
			if funcLit.Type == nil || funcLit.Type.Params == nil {
				return true
			}

			for _, field := range funcLit.Type.Params.List {
				for _, ident := range field.Names {
					if ident.Name == "main" {
						ast.Inspect(funcLit.Body, func(n ast.Node) bool {
							if callExpr, ok := n.(*ast.CallExpr); ok {
								if ident, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
									if ident.Sel.Name == "Exit" && ident.X.(*ast.Ident).Name == "os" {
										pass.Reportf(callExpr.Pos(), "вызывается os.Exit() в функции main()")
									}
								}
							}
							return true
						})
						return true
					}
				}
			}
		}
		return true
	}
}
