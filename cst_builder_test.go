package main

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/s2gatev/gocompat/cst"
)

func testCompat(
	t *testing.T,
	source string,
	expected InterfaceContext) {

	fileSet := token.NewFileSet()
	file, _ := parser.ParseFile(fileSet, "source.go", source, parser.ParseComments)

	actual := &InterfaceContext{
		Project: &cst.Project{Packages: map[string]*cst.Package{}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"MyInt": &cst.TypeDef{"MyInt", &cst.SimpleType{"int"}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"MyInt": &cst.TypeDef{"MyInt", &cst.Struct{map[string]*cst.Field{
						"A": &cst.Field{"A", &cst.SimpleType{"int"}},
						"B": &cst.Field{"B", &cst.SimpleType{"float32"}},
						"C": &cst.Field{"C", &cst.SimpleType{"string"}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"MyInt": &cst.TypeDef{"MyInt", &cst.Struct{map[string]*cst.Field{
						"A": &cst.Field{"A", &cst.SimpleType{"int"}},
						"B": &cst.Field{"B", &cst.Struct{map[string]*cst.Field{
							"C": &cst.Field{"C", &cst.SimpleType{"float32"}},
							"D": &cst.Field{"D", &cst.SimpleType{"string"}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"NameLength": &cst.Func{"NameLength",
						nil,
						&cst.Params{[]cst.Type{
							&cst.SimpleType{"string"}}},
						&cst.Results{[]cst.Type{
							&cst.SimpleType{"int"}}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"Something": &cst.Func{"Something",
						nil,
						&cst.Params{[]cst.Type{
							&cst.SimpleType{"string"},
							&cst.SimpleType{"string"},
							&cst.SimpleType{"...int"}}},
						&cst.Results{[]cst.Type{
							&cst.SimpleType{"int"},
							&cst.SimpleType{"bool"}}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"Something": &cst.Func{"Something",
						nil,
						&cst.Params{[]cst.Type{
							&cst.SimpleType{"string"},
							&cst.SimpleType{"string"},
							&cst.SimpleType{"...int"}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"Something": &cst.Func{"Something",
						nil,
						nil,
						&cst.Results{[]cst.Type{
							&cst.SimpleType{"int"}}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"A": &cst.Var{"A", &cst.SimpleType{"int"}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"A": &cst.Var{"A", &cst.SimpleType{"int"}},
					"B": &cst.Var{"B", &cst.SimpleType{"int"}},
					"D": &cst.Var{"D", &cst.SimpleType{"int"}},
					"S": &cst.Var{"S", &cst.SimpleType{"string"}},
					"F": &cst.Var{"F", &cst.SimpleType{"string"}},
					"G": &cst.Var{"G", &cst.SimpleType{"int"}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"A": &cst.Var{"A", &cst.SimpleType{"int"}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"A": &cst.Var{"A", &cst.SimpleType{"int"}},
					"B": &cst.Var{"B", &cst.SimpleType{"int"}},
					"D": &cst.Var{"D", &cst.SimpleType{"int"}},
					"S": &cst.Var{"S", &cst.SimpleType{"string"}},
					"F": &cst.Var{"F", &cst.SimpleType{"string"}},
					"G": &cst.Var{"G", &cst.SimpleType{"int"}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"MyStr": &cst.TypeDef{"MyStr", &cst.Struct{map[string]*cst.Field{}}},
					"Something": &cst.Func{"Something",
						&cst.Recievers{[]cst.Type{
							&cst.SimpleType{"MyStr"}}},
						&cst.Params{[]cst.Type{
							&cst.SimpleType{"int"}}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"MyStr": &cst.TypeDef{"MyStr", &cst.Struct{map[string]*cst.Field{}}},
					"Something": &cst.Func{"Something",
						&cst.Recievers{[]cst.Type{
							&cst.SimpleType{"*MyStr"}}},
						&cst.Params{[]cst.Type{
							&cst.SimpleType{"int"}}},
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
		Project: &cst.Project{
			Packages: map[string]*cst.Package{
				"p": &cst.Package{"p", map[string]cst.Node{
					"InterStringer": &cst.TypeDef{"InterStringer", &cst.Interface{map[string]*cst.Func{
						"String": &cst.Func{"String",
							nil,
							nil,
							&cst.Results{[]cst.Type{
								&cst.SimpleType{"string"}}},
						},
						"Int": &cst.Func{"Int",
							nil,
							&cst.Params{[]cst.Type{
								&cst.SimpleType{"float64"}}},
							&cst.Results{[]cst.Type{
								&cst.SimpleType{"int"}}},
						},
					}}},
				}},
			},
		},
	}

	testCompat(t, source, expected)
}
