package cst

// Interface represents an interface type node.
type Interface struct {
	Funcs map[string]*Func
}

func (older *Interface) Compare(n Node) bool {
	if newer, ok := n.(*Interface); ok {
		for name, sOlder := range older.Funcs {
			if sNewer, ok := newer.Funcs[name]; ok {
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
