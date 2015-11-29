package main

import "errors"

func compareSymbols(a, b []*Symbol) error {
	if len(a) != len(b) {
		return errors.New("Different number of symbols.")
	}

	for idx, sA := range a {
		sB := b[idx]

		if sA.Name != sB.Name {
			return errors.New("Different symbol name.")
		}

		if err := compareSymbols(sA.Types, sB.Types); err != nil {
			return err
		}
	}

	return nil
}

func Compare(a, b map[string]*Package) error {
	if len(a) != len(b) {
		return errors.New("Different number of packages.")
	}

	for pName, pA := range a {
		pB := b[pName]

		if pA.Name != pB.Name {
			return errors.New("Different package name.")
		}

		if err := compareSymbols(pA.Exported, pB.Exported); err != nil {
			return err
		}
	}

	return nil
}
