package retag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

type maker struct{}

func (m maker) MakeTag(t reflect.Type, fieldIndex int) reflect.StructTag {
	if t.Field(fieldIndex).Name[0] != 'X' {
		return `json:"-"`
	}
	return ""
}

type X struct {
	A int
}

type Structure struct {
	h1     string
	Int    int
	Str    string
	Bool   bool
	h2     int
	Xfield X
	X
}

func BenchmarkConvert(b *testing.B) {
	b.StopTimer()
	p := new(Structure)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Convert(p, maker{})
	}
}

// func TestX(test *testing.T) {
// 	// reflect.StructOf([]reflect.StructField{
// 	// 	{Name: "a", Type: reflect.TypeOf("")},
// 	// 	{Name: "b", Type: reflect.TypeOf("")},
// 	// })
// }

type VoidStruct struct{}

type FlatStruct struct {
	private int
	Omit    int
	Xport   int
}

type AnonymousVoidStruct struct {
	Xport int
	Xvoid VoidStruct
}

type Struct struct {
	Xport1  int
	Xport2  FlatStruct
	Omit    int
	private int
}

type PtrStruct struct {
	Xport1  int
	Xport2  *FlatStruct
	Omit    int
	private int
}

var mapTestCases = []MapTestCase{
	{"VoidStruct", maker{}, new(VoidStruct), `{}`},
	{"FlatStruct", maker{}, new(FlatStruct), `{"Xport":0}`},
	// {"AnonymousVoidStruct", maker{}, new(AnonymousVoidStruct), `{"Xport":0}`}, // strange error
	{"Struct", maker{}, new(Struct), `{"Xport1":0,"Xport2":{"Xport":0}}`},
	{"PtrStruct", maker{}, new(PtrStruct), `{"Xport1":0,"Xport2":null}`},
}

type MapTestCase struct {
	Name   string
	Maker  TagMaker
	Source interface{}
	Result string
}

func (c *MapTestCase) Run(test *testing.T) {
	result := Convert(c.Source, c.Maker)
	b, err := json.Marshal(result)
	if err != nil {
		test.Fatal("Unable to marshal result into json: ", err)
	}
	marshalled := string(b)
	if marshalled != c.Result {
		test.Errorf("Expect `%s` but got `%s`", c.Result, marshalled)
	}
}

func TestConvert(test *testing.T) {
	for _, testCase := range mapTestCases {
		test.Run(testCase.Name, testCase.Run)
	}
}

type VoidFirst struct {
	V struct{}
	A int32
}

type VoidLast struct {
	A int32
	V struct{}
}

type VoidMiddle struct {
	A int32
	V struct{}
	B int32
}

// https://play.golang.org/p/Kr9fk08S36
func _TestSizeOf(test *testing.T) {
	voidFieldDescr := reflect.StructField{
		Name: "Void",
		Type: reflect.TypeOf(struct{}{}),
	}
	aFieldDescr := reflect.StructField{Name: "A", Type: reflect.TypeOf(int32(0))}
	bFieldDescr := reflect.StructField{Name: "B", Type: reflect.TypeOf(int32(0))}

	list := []struct {
		name string
		t    reflect.Type
	}{
		{"Compiler-generated, the first field is a void struct",
			reflect.TypeOf(VoidFirst{})},
		{"Compiler-generated, the last field is a void struct",
			reflect.TypeOf(VoidLast{})},
		{"Compiler-generated, middle field is a void struct",
			reflect.TypeOf(VoidMiddle{})},
		{"Runtime-generated, the first field is a void struct",
			reflect.StructOf([]reflect.StructField{voidFieldDescr, aFieldDescr})},
		{"Runtime-generated, the last field is a void struct",
			reflect.StructOf([]reflect.StructField{aFieldDescr, voidFieldDescr})},
		{"Runtime-generated, middle field is a void struct",
			reflect.StructOf([]reflect.StructField{aFieldDescr, voidFieldDescr, bFieldDescr})},
	}

	for _, item := range list {
		fmt.Println(item.name)
		fmt.Println("Type name:", item.t.Name())
		fmt.Println("Type size:", item.t.Size())
		fmt.Println("Fields:")
		for i := 0; i < item.t.NumField(); i++ {
			field := item.t.Field(i)
			fmt.Printf("  %#v\n", field)
		}
		fmt.Println()
	}

	fmt.Println("Size of VoidFirst:", unsafe.Sizeof(VoidFirst{}))
	fmt.Println("Size of VoidLast:", unsafe.Sizeof(VoidLast{}))
	fmt.Println("Size of VoidMiddle:", unsafe.Sizeof(VoidMiddle{}))

	reflect.StructOf([]reflect.StructField{
		{Name: "Void", Type: reflect.TypeOf(struct{}{})},
		{Name: "A", Type: reflect.TypeOf(int32(0))},
	})
}
