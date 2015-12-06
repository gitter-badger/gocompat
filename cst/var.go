package cst

// Var represents a variable definition node.
type Var struct {
	Name string
	Type Type
}

func (older *Var) Compare(n Node) bool {
	if newer, ok := n.(*Var); ok {
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
