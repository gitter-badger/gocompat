package main

// Symbol represents a basic program component.
type Symbol struct {
	Name  string
	Nodes []Node
}

func compareSymbolNames(oldName, newName string) bool {
	if oldName == newName {
		return true
	}

	if "..."+oldName == newName {
		return true
	}

	return false
}

func (older *Symbol) Compare(n Node) bool {
	if newer, ok := n.(*Symbol); ok {
		if ok := compareSymbolNames(older.Name, newer.Name); !ok {
			return false
		}

		if len(older.Nodes) != len(newer.Nodes) {
			return false
		}

		for index, sOlder := range older.Nodes {
			sNewer := newer.Nodes[index]

			if !sOlder.Compare(sNewer) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}
