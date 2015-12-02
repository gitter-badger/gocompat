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

		if err := compareSymbols(sA.Symbols, sB.Symbols); err != nil {
			return err
		}
	}

	return nil
}

func compareDefinitions(a, b map[string]*Symbol) error {
	if len(a) != len(b) {
		return errors.New("Different number of definitions.")
	}

	for dName, dA := range a {
		dB := b[dName]

		if dA.Name != dB.Name {
			return errors.New("Different definition name.")
		}

		if err := compareSymbols(dA.Symbols, dB.Symbols); err != nil {
			return err
		}
	}

	return nil
}

func ComparePackages(a, b map[string]*Package) error {
	if len(a) != len(b) {
		return errors.New("Different number of packages.")
	}

	for pName, pA := range a {
		pB := b[pName]

		if pA.Name != pB.Name {
			return errors.New("Different package name.")
		}

		if err := compareDefinitions(pA.Symbols, pB.Symbols); err != nil {
			return err
		}
	}

	return nil
}
