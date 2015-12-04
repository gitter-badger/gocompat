package main

import (
	"go/ast"
	"go/token"
	"unicode"
)

func Sym(name string, nodes ...Node) *Symbol {
	return &Symbol{name, nodes}
}

// InterfaceContext is passed to the AST visitor in order to keep track the symbols
// part of the program interface.
type InterfaceContext struct {
	CurrentPackage *Package
	Application    *Application
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
func extractSymbols(expr ast.Node) []*Symbol {
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
				nodes := []Node{}
				for _, s := range extractSymbols(f.Type) {
					nodes = append(nodes, s)
				}
				symbols = append(symbols, Sym(n.Name, nodes...))
			}
		}
	case *ast.InterfaceType:
		for _, m := range t.Methods.List {
			for _, n := range m.Names {
				nodes := []Node{}
				for _, s := range extractSymbols(m.Type) {
					nodes = append(nodes, s)
				}
				symbols = append(symbols, Sym(n.Name, nodes...))
			}
		}
	case *ast.StarExpr:
		symbols = extractSymbols(t.X)
		for index, _ := range symbols {
			symbols[index].Name = "*" + symbols[index].Name
		}
	case *ast.FuncDecl:
		if t.Recv != nil {
			var recvSymbols []Node
			for _, f := range t.Recv.List {
				for _, _ = range f.Names {
					nodes := []Node{}
					for _, s := range extractSymbols(f.Type) {
						nodes = append(nodes, s)
					}
					recvSymbols = append(recvSymbols, nodes...)
				}
			}
			symbols = append(symbols, Sym("recv", recvSymbols...))
		}
		symbols = append(symbols, extractSymbols(t.Type)...)
	case *ast.FuncType:
		var paramSymbols []Node
		for _, f := range t.Params.List {
			for _, _ = range f.Names {
				nodes := []Node{}
				for _, s := range extractSymbols(f.Type) {
					nodes = append(nodes, s)
				}
				paramSymbols = append(paramSymbols, nodes...)
			}
			if f.Names == nil {
				nodes := []Node{}
				for _, s := range extractSymbols(f.Type) {
					nodes = append(nodes, s)
				}
				paramSymbols = append(paramSymbols, nodes...)
			}
		}
		symbols = append(symbols, Sym("params", paramSymbols...))

		var resultSymbols []Node
		if t.Results != nil {
			for _, f := range t.Results.List {
				nodes := []Node{}
				for _, s := range extractSymbols(f.Type) {
					nodes = append(nodes, s)
				}
				resultSymbols = append(resultSymbols, nodes...)
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

		if _, ok := context.Application.Packages[packageName]; !ok {
			context.Application.Packages[packageName] =
				&Package{packageName, map[string]Node{}}
		}
		context.CurrentPackage, _ = context.Application.Packages[packageName]
	}
}

func handleTypeSpec(node ast.Node, context interface{}) {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		context, _ := context.(*InterfaceContext)
		current := context.CurrentPackage

		if isExported(typeSpec.Name.Name) {
			symbols := extractSymbols(typeSpec.Type)
			if _, ok := typeSpec.Type.(*ast.StructType); ok {
				ms := map[string]Node{}
				for _, s := range symbols {
					ms[s.Name] = s
				}
				node := &Struct{Name: typeSpec.Name.Name, Fields: ms}
				current.Nodes[node.Name] = Sym("type", node)
			} else {
				nodes := []Node{}
				for _, s := range symbols {
					nodes = append(nodes, s)
				}
				symbol := &Symbol{typeSpec.Name.Name, nodes}
				current.Nodes[symbol.Name] = Sym("type", symbol)
			}

		}
	}
}

func handleFuncDecl(node ast.Node, context interface{}) {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		context, _ := context.(*InterfaceContext)
		current := context.CurrentPackage

		if isExported(funcDecl.Name.Name) {
			nodes := []Node{}
			for _, s := range extractSymbols(funcDecl) {
				nodes = append(nodes, s)
			}
			symbol := &Symbol{funcDecl.Name.Name, nodes}
			if funcDecl.Recv != nil {
				current.Nodes[symbol.Name] = Sym("method", symbol)
			} else {
				current.Nodes[symbol.Name] = Sym("func", symbol)
			}
		}
	}
}

func handleSpec(spec ast.Node, context interface{}) {
	if valueSpec, ok := spec.(*ast.ValueSpec); ok {
		context, _ := context.(*InterfaceContext)
		current := context.CurrentPackage

		if valueSpec.Type != nil {
			typeSymbols := extractSymbols(valueSpec.Type)
			nodes := []Node{}
			for _, s := range typeSymbols {
				nodes = append(nodes, s)
			}
			for _, name := range valueSpec.Names {
				symbol := &Symbol{name.Name, nodes}
				if isExported(symbol.Name) {
					current.Nodes[symbol.Name] = Sym("var", symbol)
				}
			}
		} else {
			for index, name := range valueSpec.Names {
				typeSymbols := extractSymbols(valueSpec.Values[index])
				nodes := []Node{}
				for _, s := range typeSymbols {
					nodes = append(nodes, s)
				}
				symbol := &Symbol{name.Name, nodes}
				if isExported(symbol.Name) {
					current.Nodes[symbol.Name] = Sym("var", symbol)
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
