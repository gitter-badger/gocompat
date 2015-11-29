package main

import (
	"go/ast"
	"go/token"
)

type Symbol struct {
	Name  string
	Types []*Symbol
}

type Package struct {
	Name     string
	Exported []*Symbol
}

type CompatContext struct {
	CurrentPackage *Package
	Packages       map[string]*Package
}

func handlePackage(node ast.Node, context *CompatContext) {
	if file, ok := node.(*ast.File); ok {
		packageName := file.Name.Name

		if _, ok := context.Packages[packageName]; !ok {
			context.Packages[packageName] = &Package{Name: packageName}
		}
		context.CurrentPackage, _ = context.Packages[packageName]
	}
}

func handleTypeSpec(node ast.Node, context *CompatContext) {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		current := context.CurrentPackage

		symbol := &Symbol{Name: typeSpec.Name.Name}
		current.Exported = append(current.Exported, symbol)
		symbol.Types = extractSymbols(typeSpec.Type)
	}
}

func handleFuncDecl(node ast.Node, context *CompatContext) {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		current := context.CurrentPackage

		symbol := &Symbol{Name: funcDecl.Name.Name}
		current.Exported = append(current.Exported, symbol)
		symbol.Types = extractSymbols(funcDecl.Type)
	}
}

func extractSymbols(expr ast.Expr) []*Symbol {
	switch t := expr.(type) {
	case *ast.Ident:
		return []*Symbol{&Symbol{Name: t.Name}}
	case *ast.StructType:
		types := []*Symbol{}
		for _, f := range t.Fields.List {
			for _, n := range f.Names {
				types = append(types, &Symbol{n.Name, extractSymbols(f.Type)})
			}
		}
		return types
	case *ast.FuncType:
		types := []*Symbol{}
		for _, f := range t.Params.List {
			for _, _ = range f.Names {
				types = append(types, extractSymbols(f.Type)...)
			}
		}
		for _, f := range t.Results.List {
			types = append(types, extractSymbols(f.Type)...)
		}
		return types
	default:
		return []*Symbol{}
	}
}

func ProcessFile(
	fileSet *token.FileSet,
	file *ast.File,
	context *CompatContext) {
	visitor := NewVisitor(fileSet, file, context)
	visitor.Handle(handlePackage)
	visitor.Handle(handleTypeSpec)
	visitor.Handle(handleFuncDecl)
	ast.Walk(visitor, file)
}
