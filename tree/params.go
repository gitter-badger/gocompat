package tree

type Params struct {
	Types []Type
}

func (older *Params) Compare(n Node) bool {
	if newer, ok := n.(*Params); ok {
		if len(older.Types) != len(newer.Types) {
			return false
		}

		for i, oType := range older.Types {
			nType := newer.Types[i]
			if !oType.Compare(nType) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}
