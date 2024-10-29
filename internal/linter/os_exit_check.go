package linter

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const reportMessage = "direct function os.Exit call is prohibited"

// OsExitCheck конфигурация чекера
var OsExitCheck = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "запрещает использовать прямой вызов os.Exit в функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	for _, f := range pass.Files {
		if f.Name.Name != "main" {
			continue
		}
		ast.Inspect(f, func(n ast.Node) bool { // проверяем, какой конкретный тип лежит в узле
			switch x := n.(type) {
			case *ast.CallExpr:
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					if ident, ok := s.X.(*ast.Ident); ok {
						if ident.Name == "os" && s.Sel.Name == "Exit" {
							pass.Reportf(n.Pos(), reportMessage)
						}
					}
				}
			case *ast.FuncDecl:
				if x.Name.Name != "main" {
					return false
				}
			}

			return true
		})
	}

	return nil, nil
}
