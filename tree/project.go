package tree

// Project represents a Go program with its constituent elements.
type Project struct {
	Packages map[string]*Package
}

func (older *Project) Compare(n Node) bool {
	if newer, ok := n.(*Project); ok {
		if len(older.Packages) != len(newer.Packages) {
			return false
		}

		for name, sOlder := range older.Packages {
			if sNewer, ok := newer.Packages[name]; ok {
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
