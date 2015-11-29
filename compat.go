package main

import (
	"go/ast"
	"go/token"
)

type Type struct {
	Name string
	Type *Type
}

type Symbol struct {
	Name  string
	Types []*Type
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
		symbol.Types = handleType(typeSpec.Type)
	}
}

func handleType(expr ast.Expr) []*Type {
	switch t := expr.(type) {
	case *ast.Ident:
		return []*Type{&Type{Name: t.Name}}
	case *ast.StructType:
		types := []*Type{}
		for _, f := range t.Fields.List {
			for _, n := range f.Names {
				types = append(types, &Type{n.Name, handleType(f.Type)[0]})
			}
		}
		return types
	default:
		return []*Type{}
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
	ast.Print(fileSet, file)
}
