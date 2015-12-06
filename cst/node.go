package cst

// Node represents an element of the concrete syntax tree.
type Node interface {

	// Compare performs semantically-specific comparison with another node.
	// Return true if both nodes are equal.
	Compare(Node) bool
}
