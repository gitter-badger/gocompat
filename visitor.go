package main

import (
	"go/ast"
	"go/token"
)

// AstVisitor visits the parts of an abstract syntax tree handling them in specific ways.
type AstVisitor struct {
	FileSet  *token.FileSet
	AST      *ast.File
	Context  *CompatContext
	Handlers []func(ast.Node, *CompatContext)
}

func NewVisitor(fileSet *token.FileSet, ast *ast.File, context *CompatContext) *AstVisitor {
	return &AstVisitor{fileSet, ast, context, nil}
}

func (v *AstVisitor) Handle(handler func(ast.Node, *CompatContext)) {
	v.Handlers = append(v.Handlers, handler)
}

func (v *AstVisitor) Visit(node ast.Node) ast.Visitor {
	for _, handler := range v.Handlers {
		handler(node, v.Context)
	}
	return v
}
