package cst

// Field represents a struct field node.
type Field struct {
	Name string
	Type Type
}

func (older *Field) Compare(n Node) bool {
	if newer, ok := n.(*Field); ok {
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
