package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func testCompat(
	t *testing.T,
	source string,
	expected CompatContext) {

	fileSet := token.NewFileSet()
	file, _ := parser.ParseFile(fileSet, "source.go", source, parser.ParseComments)

	actual := &CompatContext{Packages: map[string]*Package{}}
	ProcessFile(fileSet, file, actual)

	if err := ComparePackages(expected.Packages, actual.Packages); err != nil {
		t.Error(err)
	}
}

func TestSimpleTypeDeclaration(t *testing.T) {
	source := `
package p

type MyInt int
`

	expected := CompatContext{
		Packages: map[string]*Package{
			"p": Pack("p", map[string]*Symbol{
				"MyInt": Sym("MyInt", Sym("int")),
			}),
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
		Packages: map[string]*Package{
			"p": Pack("p", map[string]*Symbol{
				"MyInt": Sym("MyInt",
					Sym("A", Sym("int")),
					Sym("B", Sym("float32")),
					Sym("C", Sym("string")),
				),
			}),
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
		Packages: map[string]*Package{
			"p": Pack("p", map[string]*Symbol{
				"MyInt": Sym("MyInt",
					Sym("A", Sym("int")),
					Sym("B",
						Sym("C", Sym("float32")),
						Sym("D", Sym("string")),
					),
				),
			}),
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
		Packages: map[string]*Package{
			"p": Pack("p", map[string]*Symbol{}),
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
		Packages: map[string]*Package{
			"p": Pack("p", map[string]*Symbol{
				"NameLength": Sym("NameLength",
					Sym("string"),
					Sym("int"),
				),
			}),
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
		Packages: map[string]*Package{
			"p": Pack("p", map[string]*Symbol{
				"Something": Sym("Something",
					Sym("string"),
					Sym("string"),
					Sym("...int"),
					Sym("int"),
					Sym("bool"),
				),
			}),
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
		Packages: map[string]*Package{
			"p": Pack("p", map[string]*Symbol{}),
		},
	}

	testCompat(t, source, expected)
}
