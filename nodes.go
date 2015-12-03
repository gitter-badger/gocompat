package main

type Node interface {
	Compare(Node) bool
}

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

type Struct struct {
	Name   string
	Fields map[string]Node
}

func (older *Struct) Compare(n Node) bool {
	if newer, ok := n.(*Struct); ok {
		if ok := compareSymbolNames(older.Name, newer.Name); !ok {
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