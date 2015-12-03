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

type Package struct {
	Name    string
	Symbols map[string]*Symbol
}

func Pack(name string, symbols map[string]*Symbol) *Package {
	return &Package{name, symbols}
}

func Sym(name string, symbols ...*Symbol) *Symbol {
	return &Symbol{Name: name, Symbols: symbols}
}

type CompatContext struct {
	CurrentPackage *Package
	Packages       map[string]*Package
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

		if _, ok := context.Packages[packageName]; !ok {
			context.Packages[packageName] = Pack(packageName, map[string]*Symbol{})
		}
		context.CurrentPackage, _ = context.Packages[packageName]
	}
}

func handleTypeSpec(node ast.Node, context *CompatContext) {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		current := context.CurrentPackage

		symbol := &Symbol{Name: typeSpec.Name.Name}
		if isExported(symbol.Name) {
			symbol.Symbols = extractSymbols(typeSpec.Type)
			current.Symbols[symbol.Name] = symbol
		}
	}
}

func handleFuncDecl(node ast.Node, context *CompatContext) {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		current := context.CurrentPackage

		symbol := &Symbol{Name: funcDecl.Name.Name}
		if isExported(symbol.Name) {
			symbol.Symbols = extractSymbols(funcDecl.Type)
			current.Symbols[symbol.Name] = symbol
		}
	}
}

func handleSpec(spec ast.Node, context *CompatContext) {
	if valueSpec, ok := spec.(*ast.ValueSpec); ok {
		current := context.CurrentPackage

		if valueSpec.Type != nil {
			typeSymbols := extractSymbols(valueSpec.Type)
			for _, name := range valueSpec.Names {
				symbol := &Symbol{Name: name.Name, Symbols: typeSymbols}
				if isExported(symbol.Name) {
					current.Symbols[symbol.Name] = symbol
				}
			}
		} else {
			for index, name := range valueSpec.Names {
				typeSymbols := extractSymbols(valueSpec.Values[index])
				symbol := &Symbol{Name: name.Name, Symbols: typeSymbols}
				if isExported(symbol.Name) {
					current.Symbols[symbol.Name] = symbol
				}
			}
		}
	}
}

func handleGenDecl(node ast.Node, context *CompatContext) {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		for _, spec := range genDecl.Specs {
			handleSpec(spec, context)
		}
	}
}

func kindToType(kind token.Token) string {
	switch kind.String() {
	case "STRING":
		return "string"
	case "INT":
		return "int"
	default:
		return ""
	}
}

func extractSymbols(expr ast.Expr) []*Symbol {
	symbols := []*Symbol{}

	switch t := expr.(type) {
	case *ast.BasicLit:
		symbols = append(symbols, Sym(kindToType(t.Kind)))
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
		if t.Results != nil {
			for _, f := range t.Results.List {
				symbols = append(symbols, extractSymbols(f.Type)...)
			}
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
	visitor.Handle(handleGenDecl)
	ast.Walk(visitor, file)
	//ast.Print(fileSet, file)
}
