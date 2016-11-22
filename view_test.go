package retag

import (
	"reflect"
	"testing"
)

var (
	testView             = NewView("json", "admin")
	viewMakeTagTestCases = []ViewMakeTagTestCase{
		{"Void",
			``, `json:"-"`},
		{"Miss",
			`view:"user"`, `json:"-"`},
		{"Hit",
			`view:"admin"`, ``},
		{"HitAny",
			`view:"*"`, ``},
		{"HitExt",
			`view:"admin" json:"Name,omitempty"`,
			`json:"Name,omitempty"`},
		{"HitInList",
			`view:"user,admin"`, ``},
		{"HitInListExt",
			`view:"user,admin" json:"Name,omitempty"`,
			`json:"Name,omitempty"`},
	}
)

type ViewMakeTagTestCase struct {
	Name   string
	Tag    string
	Result string
}

func (c *ViewMakeTagTestCase) Run(test *testing.T) {
	field := reflect.StructField{
		Name: c.Name,
		Type: reflect.TypeOf(""),
		Tag:  reflect.StructTag(c.Tag),
	}
	t := reflect.StructOf([]reflect.StructField{field})
	result := testView.MakeTag(t, 0)
	if string(result) != c.Result {
		test.Errorf("Expect `%s` but got `%s` for tag `%s`", c.Result, result, c.Tag)
	}
}

func TestView_MakeTag(test *testing.T) {
	for _, testCase := range viewMakeTagTestCases {
		test.Run(testCase.Name, testCase.Run)
	}
}

type viewTestStruct struct{}

func TestView(test *testing.T) {
	// TODO: complete test
	Convert(new(viewTestStruct), NewView("json", "admin"))
}
