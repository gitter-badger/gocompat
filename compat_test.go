package main

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/s2gatev/gocompat/tree"
)

func testCompat(
	t *testing.T,
	source string,
	expected InterfaceContext) {

	fileSet := token.NewFileSet()
	file, _ := parser.ParseFile(fileSet, "source.go", source, parser.ParseComments)

	actual := &InterfaceContext{
		Project: &tree.Project{Packages: map[string]*tree.Package{}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"MyInt": &tree.TypeDef{"MyInt", &tree.SimpleType{"int"}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"MyInt": &tree.TypeDef{"MyInt", &tree.Struct{map[string]*tree.Field{
						"A": &tree.Field{"A", &tree.SimpleType{"int"}},
						"B": &tree.Field{"B", &tree.SimpleType{"float32"}},
						"C": &tree.Field{"C", &tree.SimpleType{"string"}},
					}}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"MyInt": &tree.TypeDef{"MyInt", &tree.Struct{map[string]*tree.Field{
						"A": &tree.Field{"A", &tree.SimpleType{"int"}},
						"B": &tree.Field{"B", &tree.Struct{map[string]*tree.Field{
							"C": &tree.Field{"C", &tree.SimpleType{"float32"}},
							"D": &tree.Field{"D", &tree.SimpleType{"string"}},
						}}},
					}}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"NameLength": &tree.Func{"NameLength",
						nil,
						&tree.Params{[]tree.Type{
							&tree.SimpleType{"string"}}},
						&tree.Results{[]tree.Type{
							&tree.SimpleType{"int"}}},
					},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"Something": &tree.Func{"Something",
						nil,
						&tree.Params{[]tree.Type{
							&tree.SimpleType{"string"},
							&tree.SimpleType{"string"},
							&tree.SimpleType{"...int"}}},
						&tree.Results{[]tree.Type{
							&tree.SimpleType{"int"},
							&tree.SimpleType{"bool"}}},
					},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"Something": &tree.Func{"Something",
						nil,
						&tree.Params{[]tree.Type{
							&tree.SimpleType{"string"},
							&tree.SimpleType{"string"},
							&tree.SimpleType{"...int"}},
						},
						nil,
					},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"Something": &tree.Func{"Something",
						nil,
						nil,
						&tree.Results{[]tree.Type{
							&tree.SimpleType{"int"}}},
					}},
				},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"A": &tree.Var{"A", &tree.SimpleType{"int"}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"A": &tree.Var{"A", &tree.SimpleType{"int"}},
					"B": &tree.Var{"B", &tree.SimpleType{"int"}},
					"D": &tree.Var{"D", &tree.SimpleType{"int"}},
					"S": &tree.Var{"S", &tree.SimpleType{"string"}},
					"F": &tree.Var{"F", &tree.SimpleType{"string"}},
					"G": &tree.Var{"G", &tree.SimpleType{"int"}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"A": &tree.Var{"A", &tree.SimpleType{"int"}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"A": &tree.Var{"A", &tree.SimpleType{"int"}},
					"B": &tree.Var{"B", &tree.SimpleType{"int"}},
					"D": &tree.Var{"D", &tree.SimpleType{"int"}},
					"S": &tree.Var{"S", &tree.SimpleType{"string"}},
					"F": &tree.Var{"F", &tree.SimpleType{"string"}},
					"G": &tree.Var{"G", &tree.SimpleType{"int"}},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"MyStr": &tree.TypeDef{"MyStr", &tree.Struct{map[string]*tree.Field{}}},
					"Something": &tree.Func{"Something",
						&tree.Recievers{[]tree.Type{
							&tree.SimpleType{"MyStr"}}},
						&tree.Params{[]tree.Type{
							&tree.SimpleType{"int"}}},
						nil,
					},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"MyStr": &tree.TypeDef{"MyStr", &tree.Struct{map[string]*tree.Field{}}},
					"Something": &tree.Func{"Something",
						&tree.Recievers{[]tree.Type{
							&tree.SimpleType{"*MyStr"}}},
						&tree.Params{[]tree.Type{
							&tree.SimpleType{"int"}}},
						nil,
					},
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
		Project: &tree.Project{
			Packages: map[string]*tree.Package{
				"p": &tree.Package{"p", map[string]tree.Node{
					"InterStringer": &tree.TypeDef{"InterStringer", &tree.Interface{map[string]*tree.Func{
						"String": &tree.Func{"String",
							nil,
							nil,
							&tree.Results{[]tree.Type{
								&tree.SimpleType{"string"}}},
						},
						"Int": &tree.Func{"Int",
							nil,
							&tree.Params{[]tree.Type{
								&tree.SimpleType{"float64"}}},
							&tree.Results{[]tree.Type{
								&tree.SimpleType{"int"}}},
						},
					}}},
				}},
			},
		},
	}

	testCompat(t, source, expected)
}
