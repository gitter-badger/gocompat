package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func parse(source string) *Application {
	fileSet := token.NewFileSet()
	file, _ := parser.ParseFile(fileSet, "source.go", source, parser.ParseComments)

	context := &InterfaceContext{
		Application: &Application{Packages: map[string]*Package{}},
	}
	ProcessFile(fileSet, file, context)

	return context.Application
}

func testCompare(
	t *testing.T,
	older, newer string,
	shouldHaveError bool) {

	ok := parse(older).Compare(parse(newer))

	if ok == shouldHaveError {
		t.Error("Error in compare test.")
	}
}

func TestAddSymbolsToPackage(t *testing.T) {
	older := `
package p

type MyInt int
`

	newer := `
package p

type MyInt int
type MyFloat float64
`

	testCompare(t, older, newer, false)
}

func TestRemoveSymbolsFromPackage(t *testing.T) {
	older := `
package p

type MyInt int
type MyFloat float64
`

	newer := `
package p

type MyFloat float64
`

	testCompare(t, older, newer, true)
}

func TestChangeTypeBase(t *testing.T) {
	older := `
package p

type MyFloat float64
`

	newer := `
package p

type MyFloat int32
`

	testCompare(t, older, newer, true)
}

func TestChangeVarType(t *testing.T) {
	older := `
package p

var A = 42
`

	newer := `
package p

var A string = "answer"
`

	testCompare(t, older, newer, true)
}

func TestChangeConstType(t *testing.T) {
	older := `
package p

const A = 42
`

	newer := `
package p

const A string = "answer"
`

	testCompare(t, older, newer, true)
}

func TestFuncChangeArgType(t *testing.T) {
	older := `
package p

func Something(a int) {
}
`

	newer := `
package p

func Something(a string) {
}
`

	testCompare(t, older, newer, true)
}

func TestFuncArgToVarArgs(t *testing.T) {
	older := `
package p

func Something(a int) {
}
`

	newer := `
package p

func Something(a ...int) {
}
`

	testCompare(t, older, newer, false)
}

func TestTypeToVar(t *testing.T) {
	older := `
package p

type A int
`

	newer := `
package p

var A = 5
`

	testCompare(t, older, newer, true)
}

func TestVarToFunc(t *testing.T) {
	older := `
package p

var A int = 5
`

	newer := `
package p

func A(a int) {
}
`

	testCompare(t, older, newer, true)
}

func TestFuncToType(t *testing.T) {
	older := `
package p

func A(a int) {
}
`

	newer := `
package p

type A struct {}
`

	testCompare(t, older, newer, true)
}

func TestFuncParamToResult(t *testing.T) {
	older := `
package p

func A(a int) {
}
`

	newer := `
package p

func A() int {
	return 42
}
`

	testCompare(t, older, newer, true)
}

func TestFuncResultToParam(t *testing.T) {
	older := `
package p

func A() (string, int) {
	return "", 42
}
`

	newer := `
package p

func A(a string, b int) {
}
`

	testCompare(t, older, newer, true)
}
