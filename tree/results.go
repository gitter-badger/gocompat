package tree

type Results struct {
	Types []Type
}

func (older *Results) Compare(n Node) bool {
	if newer, ok := n.(*Results); ok {
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
