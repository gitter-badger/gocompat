package tree

// Struct represents a struct element in a Go program.
type Struct struct {
	Fields map[string]*Field
}

func (older *Struct) Compare(n Node) bool {
	if newer, ok := n.(*Struct); ok {
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
