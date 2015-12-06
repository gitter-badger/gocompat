package tree

type Recievers struct {
	Types []Type
}

func (older *Recievers) Compare(n Node) bool {
	if newer, ok := n.(*Recievers); ok {
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
