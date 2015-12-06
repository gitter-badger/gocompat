package cst

// Func represents a function definition node.
type Func struct {
	Name      string
	Recievers *Recievers
	Params    *Params
	Results   *Results
}

func (older *Func) Compare(n Node) bool {
	if newer, ok := n.(*Func); ok {

		if older.Recievers == newer.Recievers {
		} else if older.Recievers == nil && newer.Recievers != nil {
			return false
		} else if older.Recievers != nil && newer.Recievers == nil {
			return false
		} else if !older.Recievers.Compare(newer.Recievers) {
			return false
		}

		if older.Params == newer.Params {
		} else if older.Params == nil && newer.Params != nil {
			return false
		} else if older.Params != nil && newer.Params == nil {
			return false
		} else if !older.Params.Compare(newer.Params) {
			return false
		}

		if older.Results == newer.Results {
		} else if older.Results == nil && newer.Results != nil {
			return false
		} else if older.Results != nil && newer.Results == nil {
			return false
		} else if !older.Results.Compare(newer.Results) {
			return false
		}

		return true
	} else {
		return false
	}
}
