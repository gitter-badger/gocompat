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

// InterfaceContext is passed to the AST visitor in order to keep track the symbols
// part of the program interface.
type InterfaceContext struct {
	CurrentPackage *Package
	Packages       map[string]*Package
}

// isExporeted returns if a given name should be public or private.
func isExported(name string) bool {
	for _, r := range name {
		return unicode.IsUpper(r)
	}
	return false
}

// kindToType transforms Go token kind to type name.
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

// extractSymbols returns the interface-specific symbols part of an AST expression.
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
		var paramSymbols []*Symbol
		for _, f := range t.Params.List {
			for _, _ = range f.Names {
				paramSymbols = append(paramSymbols, extractSymbols(f.Type)...)
			}
		}
		symbols = append(symbols, Sym("params", paramSymbols...))

		var resultSymbols []*Symbol
		if t.Results != nil {
			for _, f := range t.Results.List {
				resultSymbols = append(resultSymbols, extractSymbols(f.Type)...)
			}
		}
		symbols = append(symbols, Sym("results", resultSymbols...))
	}

	return symbols
}

func handlePackage(node ast.Node, context interface{}) {
	if file, ok := node.(*ast.File); ok {
		context, _ := context.(*InterfaceContext)
		packageName := file.Name.Name

		if _, ok := context.Packages[packageName]; !ok {
			context.Packages[packageName] = Pack(packageName, map[string]*Symbol{})
		}
		context.CurrentPackage, _ = context.Packages[packageName]
	}
}

func handleTypeSpec(node ast.Node, context interface{}) {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		context, _ := context.(*InterfaceContext)
		current := context.CurrentPackage

		symbol := &Symbol{Name: typeSpec.Name.Name}
		if isExported(symbol.Name) {
			symbol.Symbols = extractSymbols(typeSpec.Type)
			current.Symbols[symbol.Name] = Sym("type", symbol)
		}
	}
}

func handleFuncDecl(node ast.Node, context interface{}) {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		context, _ := context.(*InterfaceContext)
		current := context.CurrentPackage

		symbol := &Symbol{Name: funcDecl.Name.Name}
		if isExported(symbol.Name) {
			symbol.Symbols = extractSymbols(funcDecl.Type)
			current.Symbols[symbol.Name] = Sym("func", symbol)
		}
	}
}

func handleSpec(spec ast.Node, context interface{}) {
	if valueSpec, ok := spec.(*ast.ValueSpec); ok {
		context, _ := context.(*InterfaceContext)
		current := context.CurrentPackage

		if valueSpec.Type != nil {
			typeSymbols := extractSymbols(valueSpec.Type)
			for _, name := range valueSpec.Names {
				symbol := &Symbol{Name: name.Name, Symbols: typeSymbols}
				if isExported(symbol.Name) {
					current.Symbols[symbol.Name] = Sym("var", symbol)
				}
			}
		} else {
			for index, name := range valueSpec.Names {
				typeSymbols := extractSymbols(valueSpec.Values[index])
				symbol := &Symbol{Name: name.Name, Symbols: typeSymbols}
				if isExported(symbol.Name) {
					current.Symbols[symbol.Name] = Sym("var", symbol)
				}
			}
		}
	}
}

func handleGenDecl(node ast.Node, context interface{}) {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		context, _ := context.(*InterfaceContext)

		for _, spec := range genDecl.Specs {
			handleSpec(spec, context)
		}
	}
}

func ProcessFile(
	fileSet *token.FileSet,
	file *ast.File,
	context *InterfaceContext) {

	visitor := &ContextPassingVisitor{FileSet: fileSet, AST: file, Context: context}
	visitor.Handle(handlePackage)
	visitor.Handle(handleTypeSpec)
	visitor.Handle(handleFuncDecl)
	visitor.Handle(handleGenDecl)

	ast.Walk(visitor, file)
}
