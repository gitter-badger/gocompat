package main

import (
	"go/ast"
	"go/token"
)

type Symbol struct {
	Name  string
	Types []string
}

type Package struct {
	Name     string
	Exported []*Symbol
}

type CompatContext struct {
	CurrentPackage *Package
	CurrentSymbol  *Symbol
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
		context.CurrentSymbol = symbol
		handleType(typeSpec.Type, context)
	}
}

func handleType(expr ast.Expr, context *CompatContext) {
	if object, ok := expr.(*ast.Ident); ok {
		current := context.CurrentSymbol

		current.Types = append(current.Types, object.Name)
	}
}

func ProcessFile(
	fileSet *token.FileSet,
	file *ast.File,
	context *CompatContext) {
	visitor := NewVisitor(fileSet, file, context)
	visitor.Handle(handlePackage)
	visitor.Handle(handleTypeSpec)
	ast.Walk(visitor, file)
}
