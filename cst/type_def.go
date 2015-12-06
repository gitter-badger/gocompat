package cst

// TypeDef represents a type definition node.
type TypeDef struct {
	Name string
	Type Type
}

func (older *TypeDef) Compare(n Node) bool {
	if newer, ok := n.(*TypeDef); ok {
		if older.Name != newer.Name {
			return false
		}

		if !older.Type.Compare(newer.Type) {
			return false
		}

		return true
	} else {
		return false
	}
}
