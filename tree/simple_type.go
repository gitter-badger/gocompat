package tree

type SimpleType struct {
	Name string
}

func (older *SimpleType) Compare(n Node) bool {
	if newer, ok := n.(*SimpleType); ok {
		return older.Name == newer.Name
	} else {
		return false
	}
}
