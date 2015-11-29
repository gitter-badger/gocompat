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

func Sym(name string, symbols ...*Symbol) *Symbol {
	return &Symbol{Name: name, Symbols: symbols}
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
	symbols := []*Symbol{}

	switch t := expr.(type) {
	case *ast.Ident:
		symbols = append(symbols, Sym(t.Name))
	case *ast.Ellipsis:
		symbols = extractSymbols(t.Elt)
		for index, _ := range symbols {
			symbols[index].Name = "..." + symbols[index].Name
		}
	case *ast.StructType:
		for _, f := range t.Fields.List {
			for _, n := range f.Names {
				symbols = append(symbols, Sym(n.Name, extractSymbols(f.Type)...))
			}
		}
	case *ast.FuncType:
		for _, f := range t.Params.List {
			for _, _ = range f.Names {
				symbols = append(symbols, extractSymbols(f.Type)...)
			}
		}
		for _, f := range t.Results.List {
			symbols = append(symbols, extractSymbols(f.Type)...)
		}
	}

	return symbols
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
