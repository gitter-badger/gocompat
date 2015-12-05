package tree

// Package represents a package element in a Go program.
type Package struct {
	Name  string
	Nodes map[string]Node
}

func (older *Package) Compare(n Node) bool {
	if newer, ok := n.(*Package); ok {
		for name, sOlder := range older.Nodes {
			if sNewer, ok := newer.Nodes[name]; ok {
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
