package main

// Node represents an element of a computer program developed in Go.
type Node interface {

	// Compare performs semantically-specific comparison with another node.
	// Return true if both nodes are equal.
	Compare(Node) bool
}
