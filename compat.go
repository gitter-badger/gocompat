package main

import (
	"go/ast"
	"go/token"
	"unicode"

	"github.com/s2gatev/gocompat/tree"
)

// InterfaceContext is passed to the AST visitor in order to keep track the symbols
// part of the program interface.
type InterfaceContext struct {
	CurrentPackage *tree.Package
	Project        *tree.Project
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

func extractTypes(node ast.Node) []tree.Type {
	types := []tree.Type{}
	switch n := node.(type) {
	case *ast.BasicLit:
		types = append(types, &tree.SimpleType{kindToType(n.Kind)})
	case *ast.Ident:
		types = append(types, &tree.SimpleType{n.Name})
	case *ast.Ellipsis:
		types = extractTypes(n.Elt)
		for _, t := range types {
			if st, ok := t.(*tree.SimpleType); ok {
				st.Name = "..." + st.Name
			}
		}
	case *ast.StarExpr:
		types = extractTypes(n.X)
		for _, t := range types {
			if st, ok := t.(*tree.SimpleType); ok {
				st.Name = "*" + st.Name
			}
		}
	case *ast.StructType:
		types = append(types, extractStruct(n))
	}
	return types
}

func extractFields(s *ast.StructType) map[string]*tree.Field {
	fields := map[string]*tree.Field{}
	for _, f := range s.Fields.List {
		for _, n := range f.Names {
			fields[n.Name] = &tree.Field{n.Name, extractTypes(f.Type)[0]}
		}
	}
	return fields
}

func extractFuncs(i *ast.InterfaceType) map[string]*tree.Func {
	funcs := map[string]*tree.Func{}
	for _, f := range i.Methods.List {
		for _, n := range f.Names {
			params, results := extractFuncTypeDefinition(f.Type.(*ast.FuncType))
			funcs[n.Name] = &tree.Func{n.Name, nil, params, results}
		}
	}
	return funcs
}

func extractStruct(s *ast.StructType) *tree.Struct {
	return &tree.Struct{Fields: extractFields(s)}
}

func extractFuncTypeDefinition(f *ast.FuncType) (*tree.Params, *tree.Results) {
	var params *tree.Params
	var results *tree.Results

	// Extract params.
	if f.Params != nil && f.Params.List != nil {
		var paramTypes []tree.Type
		for _, p := range f.Params.List {
			for _, _ = range p.Names {
				paramTypes = append(paramTypes, extractTypes(p.Type)...)
			}
			if p.Names == nil {
				paramTypes = append(paramTypes, extractTypes(p.Type)...)
			}
		}
		params = &tree.Params{paramTypes}
	}

	// Extract results.
	if f.Results != nil {
		var resultTypes []tree.Type
		for _, r := range f.Results.List {
			resultTypes = append(resultTypes, extractTypes(r.Type)...)
		}
		results = &tree.Results{resultTypes}
	}

	return params, results
}

func extractFuncDefinition(f *ast.FuncDecl) (*tree.Recievers, *tree.Params, *tree.Results) {
	var recievers *tree.Recievers
	var params *tree.Params
	var results *tree.Results

	// Extract recievers.
	if f.Recv != nil {
		var recieverTypes []tree.Type
		for _, r := range f.Recv.List {
			for _, _ = range r.Names {
				recieverTypes = append(recieverTypes, extractTypes(r.Type)...)
			}
		}
		recievers = &tree.Recievers{recieverTypes}
	}

	params, results = extractFuncTypeDefinition(f.Type)

	return recievers, params, results
}

func handlePackage(node ast.Node, context interface{}) {
	if file, ok := node.(*ast.File); ok {
		context, _ := context.(*InterfaceContext)
		packageName := file.Name.Name

		if _, ok := context.Project.Packages[packageName]; !ok {
			context.Project.Packages[packageName] =
				&tree.Package{packageName, map[string]tree.Node{}}
		}
		context.CurrentPackage, _ = context.Project.Packages[packageName]
	}
}

func handleTypeSpec(node ast.Node, context interface{}) {
	if typeSpec, ok := node.(*ast.TypeSpec); ok {
		context, _ := context.(*InterfaceContext)
		current := context.CurrentPackage

		if isExported(typeSpec.Name.Name) {
			var st tree.Type
			switch t := typeSpec.Type.(type) {
			case *ast.StructType:
				fields := extractFields(t)
				st = &tree.Struct{fields}
			case *ast.InterfaceType:
				funcs := extractFuncs(t)
				st = &tree.Interface{funcs}
			default:
				st = extractTypes(t)[0]
			}
			current.Nodes[typeSpec.Name.Name] = &tree.TypeDef{typeSpec.Name.Name, st}
		}
	}
}

func handleFuncDecl(node ast.Node, context interface{}) {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		context, _ := context.(*InterfaceContext)
		current := context.CurrentPackage

		if isExported(funcDecl.Name.Name) {
			recievers, params, results := extractFuncDefinition(funcDecl)
			current.Nodes[funcDecl.Name.Name] = &tree.Func{
				funcDecl.Name.Name, recievers, params, results}
		}
	}
}

func handleSpec(spec ast.Node, context interface{}) {
	if valueSpec, ok := spec.(*ast.ValueSpec); ok {
		context, _ := context.(*InterfaceContext)
		current := context.CurrentPackage

		if valueSpec.Type != nil {
			varTypes := extractTypes(valueSpec.Type)
			for _, name := range valueSpec.Names {
				varSpec := &tree.Var{name.Name, varTypes[0]}
				if isExported(name.Name) {
					current.Nodes[name.Name] = varSpec
				}
			}
		} else {
			for index, name := range valueSpec.Names {
				varTypes := extractTypes(valueSpec.Values[index])
				varSpec := &tree.Var{name.Name, varTypes[0]}
				if isExported(name.Name) {
					current.Nodes[name.Name] = varSpec
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
