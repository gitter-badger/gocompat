package main

import (
	"go/ast"
	"go/token"
)

// ContextPassingVisitor visits the nodes of an abstract syntax tree
// handling them in specific ways, passing a context object.
type ContextPassingVisitor struct {
	FileSet  *token.FileSet
	AST      *ast.File
	Context  interface{}
	Handlers []func(ast.Node, interface{})
}

// Handle assings a new handler to the visitor.
func (cpv *ContextPassingVisitor) Handle(handler func(ast.Node, interface{})) {
	cpv.Handlers = append(cpv.Handlers, handler)
}

// Visit traverses the ast applying visitor's handlers to each node in the order they are defined.
func (cpv *ContextPassingVisitor) Visit(node ast.Node) ast.Visitor {
	for _, handler := range cpv.Handlers {
		handler(node, cpv.Context)
	}
	return cpv
}
