package main

import "errors"

func packageMissingErr(packageName string) error {
	return errors.New("Package " + packageName + " is missing.")
}

func corruptedPackageErr(packageName string) error {
	return errors.New("Package " + packageName + " is corrupted.")
}

func definitionMissingErr(definitionName string) error {
	return errors.New("Definition of " + definitionName + " is missing.")
}

func corruptedDefinitionErr(definitionName string) error {
	return errors.New("Definition of " + definitionName + " is corrupted.")
}

func corruptedSymbolErr() error {
	return errors.New("Corrupted symbol.")
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

func compareSymbols(older, newer []*Symbol) error {
	if len(older) != len(newer) {
		return corruptedSymbolErr()
	}

	for index, sOlder := range older {
		sNewer := newer[index]

		if !compareSymbolNames(sOlder.Name, sNewer.Name) {
			return corruptedSymbolErr()
		}

		if err := compareSymbols(sOlder.Symbols, sNewer.Symbols); err != nil {
			return err
		}
	}

	return nil
}

func compareDefinitions(older, newer map[string]*Symbol) error {
	for dName, dOlder := range older {
		dNewer, definitionExists := newer[dName]

		if !definitionExists {
			return definitionMissingErr(dName)
		}

		if dNewer.Name != dOlder.Name {
			return corruptedDefinitionErr(dName)
		}

		if err := compareSymbols(dOlder.Symbols, dNewer.Symbols); err != nil {
			return err
		}
	}

	return nil
}

func ComparePackages(older, newer map[string]*Package) error {
	for pName, pOlder := range older {
		pNewer, packageExists := newer[pName]

		if !packageExists {
			return packageMissingErr(pName)
		}

		if pNewer.Name != pOlder.Name {
			return corruptedPackageErr(pName)
		}

		if err := compareDefinitions(pOlder.Symbols, pNewer.Symbols); err != nil {
			return err
		}
	}

	return nil
}
