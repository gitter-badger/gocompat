package main

type Node interface {
	Compare(Node) bool
}

type Symbol struct {
	Name    string
	Symbols []*Symbol
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

		if len(older.Symbols) != len(newer.Symbols) {
			return false
		}

		for index, sOlder := range older.Symbols {
			sNewer := newer.Symbols[index]

			if !sOlder.Compare(sNewer) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

type Package struct {
	Name    string
	Symbols map[string]*Symbol
}

func (older *Package) Compare(n Node) bool {
	if newer, ok := n.(*Package); ok {
		for name, sOlder := range older.Symbols {
			if sNewer, ok := newer.Symbols[name]; ok {
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

type Application struct {
	Packages map[string]*Package
}

func (older *Application) Compare(n Node) bool {
	if newer, ok := n.(*Application); ok {
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
