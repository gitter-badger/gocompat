package main

import (
	"go/ast"
	"go/token"
	"unicode"
)

type Symbol struct {
	Name    string
	Symbols []*Symbol
}

type CompatContext struct {
	CurrentSymbol *Symbol
	Symbols       map[string]*Symbol
}

func isExported(name string) bool {
	for _, r := range name {
		return unicode.IsUpper(r)
	}
	return false
}

func handlePackage(node ast.Node, context *CompatContext) {
	if file, ok := node.(*ast.File); ok {
		packageName := file.Name.Name

		if _, ok := context.Symbols[packageName]; !ok {
			context.Symbols[packageName] = &Symbol{Name: packageName}
		}
		context.CurrentSymbol, _ = context.Symbols[packageName]
	}
}

func handleTypeSpec(node ast.Node, context *CompatContext) {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		current := context.CurrentSymbol

		symbol := &Symbol{Name: typeSpec.Name.Name}
		if isExported(symbol.Name) {
			symbol.Symbols = extractSymbols(typeSpec.Type)
			current.Symbols = append(current.Symbols, symbol)
		}
	}
}

func handleFuncDecl(node ast.Node, context *CompatContext) {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		current := context.CurrentSymbol

		symbol := &Symbol{Name: funcDecl.Name.Name}
		if isExported(symbol.Name) {
			symbol.Symbols = extractSymbols(funcDecl.Type)
			current.Symbols = append(current.Symbols, symbol)
		}
	}
}

func extractSymbols(expr ast.Expr) []*Symbol {
	switch t := expr.(type) {
	case *ast.Ident:
		return []*Symbol{&Symbol{Name: t.Name}}
	case *ast.Ellipsis:
		types := extractSymbols(t.Elt)
		for index, _ := range types {
			types[index].Name = "..." + types[index].Name
		}
		return types
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
