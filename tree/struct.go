package tree

// Struct represents a struct element in a Go program.
type Struct struct {
	Name   string
	Fields map[string]Node
}

func (older *Struct) Compare(n Node) bool {
	if newer, ok := n.(*Struct); ok {
		if older.Name != newer.Name {
			return false
		}

		for name, sOlder := range older.Fields {
			if sNewer, ok := newer.Fields[name]; ok {
				if !sOlder.Compare(sNewer) {
					return false
				}
			} else {
				return false
			}
		}
		return true
	} else {
		return false
	}
}
