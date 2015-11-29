package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func sym(value string, subs ...*Symbol) *Symbol {
	return &Symbol{value, subs}
}

func testCompat(
	t *testing.T,
	source string,
	expected CompatContext) {

	fileSet := token.NewFileSet()
	file, _ := parser.ParseFile(fileSet, "source.go", source, parser.ParseComments)

	actual := &CompatContext{Symbols: map[string]*Symbol{}}
	ProcessFile(fileSet, file, actual)

	if err := Compare(expected.Symbols, actual.Symbols); err != nil {
		t.Error(err)
	}
}

func TestSimpleTypeDeclaration(t *testing.T) {
	source := `
package p

type MyInt int
`

	expected := CompatContext{
		Symbols: map[string]*Symbol{
			"p": sym("p",
				sym("MyInt",
					sym("int"),
				),
			),
		},
	}

	testCompat(t, source, expected)
}

func TestStructTypeDeclaration(t *testing.T) {
	source := `
package p

type MyInt struct {
	A	int
	B	float32
	C	string
}
`

	expected := CompatContext{
		Symbols: map[string]*Symbol{
			"p": sym("p",
				sym("MyInt",
					sym("A", sym("int")),
					sym("B", sym("float32")),
					sym("C", sym("string")),
				),
			),
		},
	}

	testCompat(t, source, expected)
}

func TestNestedStructTypeDeclaration(t *testing.T) {
	source := `
package p

type MyInt struct {
	A	int
	B	struct {
		C	float32
		D	string
	}
}
`

	expected := CompatContext{
		Symbols: map[string]*Symbol{
			"p": sym("p",
				sym("MyInt",
					sym("A", sym("int")),
					sym("B",
						sym("C", sym("float32")),
						sym("D", sym("string")),
					),
				),
			),
		},
	}

	testCompat(t, source, expected)
}

func TestNotExportedTypeDeclaration(t *testing.T) {
	source := `
package p

type myInt int
`

	expected := CompatContext{
		Symbols: map[string]*Symbol{
			"p": sym("p"),
		},
	}

	testCompat(t, source, expected)
}

func TestSimpleFuncDeclaration(t *testing.T) {
	source := `
package p

func NameLength(name string) int {
	return len(name)
}
`

	expected := CompatContext{
		Symbols: map[string]*Symbol{
			"p": sym("p",
				sym("NameLength",
					sym("string"),
					sym("int"),
				),
			),
		},
	}

	testCompat(t, source, expected)
}

func TestComplexFuncDeclaration(t *testing.T) {
	source := `
package p

func Something(a, b string, options ...int) (int, bool) {
	return 42, true
}
`

	expected := CompatContext{
		Symbols: map[string]*Symbol{
			"p": sym("p",
				sym("Something",
					sym("string"),
					sym("string"),
					sym("...int"),
					sym("int"),
					sym("bool"),
				),
			),
		},
	}

	testCompat(t, source, expected)
}

func TestNotExportedFunc(t *testing.T) {
	source := `
package p

func something(a, b string, options ...int) (int, bool) {
	return 42, true
}
`

	expected := CompatContext{
		Symbols: map[string]*Symbol{
			"p": sym("p"),
		},
	}

	testCompat(t, source, expected)
}
