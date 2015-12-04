package main

import (
	"go/parser"
	"go/token"
	"testing"
)

func testCompat(
	t *testing.T,
	source string,
	expected InterfaceContext) {

	fileSet := token.NewFileSet()
	file, _ := parser.ParseFile(fileSet, "source.go", source, parser.ParseComments)

	actual := &InterfaceContext{
		Project: &Project{Packages: map[string]*Package{}},
	}
	ProcessFile(fileSet, file, actual)

	if ok := expected.Project.Compare(actual.Project); !ok {
		t.Error("Error in compat test.")
	}
}

func TestSimpleTypeDeclaration(t *testing.T) {
	source := `
package p

type MyInt int
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"MyInt": Sym("type", Sym("MyInt", Sym("int"))),
				}},
			},
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

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"MyInt": Sym("type",
						&Struct{"MyInt", map[string]Node{
							"A": Sym("A", Sym("int")),
							"B": Sym("B", Sym("float32")),
							"C": Sym("C", Sym("string")),
						}}),
				}},
			},
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

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"MyInt": Sym("type",
						&Struct{"MyInt", map[string]Node{
							"A": Sym("A", Sym("int")),
							"B": Sym("B",
								Sym("C", Sym("float32")),
								Sym("D", Sym("string")),
							),
						}}),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestNotExportedTypeDeclaration(t *testing.T) {
	source := `
package p

type myInt int
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{}},
			},
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

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"NameLength": Sym("func", Sym("NameLength",
						Sym("params",
							Sym("string")),
						Sym("results",
							Sym("int")),
					)),
				}},
			},
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

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"Something": Sym("func", Sym("Something",
						Sym("params",
							Sym("string"),
							Sym("string"),
							Sym("...int")),
						Sym("results",
							Sym("int"),
							Sym("bool")),
					)),
				}},
			},
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

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestFuncWithoutResults(t *testing.T) {
	source := `
package p

func Something(a, b string, options ...int) {
}
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"Something": Sym("func", Sym("Something",
						Sym("params",
							Sym("string"),
							Sym("string"),
							Sym("...int")),
						Sym("results"))),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestFuncWithoutParams(t *testing.T) {
	source := `
package p

func Something() int {
	return 42
}
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"Something": Sym("func", Sym("Something",
						Sym("params"),
						Sym("results",
							Sym("int")))),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestExportedVar(t *testing.T) {
	source := `
package p

var A int = 5
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"A": Sym("var", Sym("A", Sym("int"))),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestNotExportedVar(t *testing.T) {
	source := `
package p

var a int = 5
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestComplexVar(t *testing.T) {
	source := `
package p

var A, B, c, D int = 5
var S string = "something"
var F, G = "answer", 42
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"A": Sym("var", Sym("A", Sym("int"))),
					"B": Sym("var", Sym("B", Sym("int"))),
					"D": Sym("var", Sym("D", Sym("int"))),
					"S": Sym("var", Sym("S", Sym("string"))),
					"F": Sym("var", Sym("F", Sym("string"))),
					"G": Sym("var", Sym("G", Sym("int"))),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestExportedConst(t *testing.T) {
	source := `
package p

const A int = 5
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"A": Sym("var", Sym("A", Sym("int"))),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestNotExportedConst(t *testing.T) {
	source := `
package p

const a int = 5
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestComplexConst(t *testing.T) {
	source := `
package p

const A, B, c, D int = 5
const S string = "something"
const F, G = "answer", 42
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"A": Sym("var", Sym("A", Sym("int"))),
					"B": Sym("var", Sym("B", Sym("int"))),
					"D": Sym("var", Sym("D", Sym("int"))),
					"S": Sym("var", Sym("S", Sym("string"))),
					"F": Sym("var", Sym("F", Sym("string"))),
					"G": Sym("var", Sym("G", Sym("int"))),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestStructMethod(t *testing.T) {
	source := `
package p

type MyStr struct {}

func (ms MyStr) Something(a int) {
}
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"MyStr": Sym("type", &Struct{"MyStr", map[string]Node{}}),
					"Something": Sym("method",
						Sym("Something",
							Sym("recv",
								Sym("MyStr")),
							Sym("params",
								Sym("int")),
							Sym("results"))),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestStructPointerMethod(t *testing.T) {
	source := `
package p

type MyStr struct {}

func (ms *MyStr) Something(a int) {
}
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"MyStr": Sym("type", &Struct{"MyStr", map[string]Node{}}),
					"Something": Sym("method",
						Sym("Something",
							Sym("recv",
								Sym("*MyStr")),
							Sym("params",
								Sym("int")),
							Sym("results"))),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}

func TestInterface(t *testing.T) {
	source := `
package p

type InterStringer interface {
	String() string
	Int(float64) int
}
`

	expected := InterfaceContext{
		Project: &Project{
			Packages: map[string]*Package{
				"p": &Package{"p", map[string]Node{
					"InterStringer": Sym("type",
						Sym("InterStringer",
							Sym("String",
								Sym("params"),
								Sym("results",
									Sym("string"))),
							Sym("Int",
								Sym("params",
									Sym("float64")),
								Sym("results",
									Sym("int"))))),
				}},
			},
		},
	}

	testCompat(t, source, expected)
}
