package main

import (
	"encoding/json"

	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

type Param struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}

type Result struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type Func struct {
	Name string `json:"name"`

	Params []Param `json:"params,omitempty"`
	Results []Result `json:"results,omitempty"`
}

var (
	Funcs map[string]*Func
)

func VisitFuncs(node ast.Node) bool {
	if node == nil {
		return false
	}
		
	switch n := node.(type) {
	case *ast.File, *ast.Ident, *ast.FuncType:
		return true

	case *ast.FuncDecl:
		newFunc := &Func{ Name: n.Name.Name }

		Funcs[n.Name.Name] = newFunc

		if n.Type.Params != nil {
			for _, fieldParam := range n.Type.Params.List {
				param := Param{
					Name: fieldParam.Names[0].Name, 
					Type: types.ExprString(fieldParam.Type)}

				newFunc.Params = append(newFunc.Params, param)
			}
		} 

		if n.Type.Results != nil {
			for _, fieldResult := range n.Type.Results.List {
				result := Result{
					Type: types.ExprString(fieldResult.Type)}

				if len(fieldResult.Names) > 0 {
					result.Name = fieldResult.Names[0].Name
				}

				newFunc.Results = append(newFunc.Results, result)
			}
		}
	}

	return false
}

func ParseSymbols(path string) (string, error) {
	Funcs = make(map[string]*Func)

	fset := token.NewFileSet()

	parsed, e := parser.ParseDir(fset, path, nil, 0)
	if e != nil {
		return "", e
	}

	for _, pkg := range parsed {
		ast.PackageExports(pkg)
	}

	for _, astpkg := range parsed {
		for _, f := range astpkg.Files {
			ast.Inspect(f, VisitFuncs)
		}
	}

	res, e := json.Marshal(&Funcs)
	if e != nil {
		return "", e
	}

	return string(res), nil
}
