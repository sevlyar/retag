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

type VoidStruct struct{}

type FlatStruct struct {
	private string
	Omit    int
	Xport   int
}

type AnonymousVoidStruct struct {
	Xport int
	Xvoid VoidStruct
}

type Struct struct {
	Xport1 int
	Xport2 FlatStruct
}

type PtrStruct struct {
	Xport1 int
	Xport2 *FlatStruct
}

type SliceStruct struct {
	XportSlice []FlatStruct
	XportArray [2]FlatStruct
}

type MapStruct struct {
	XportMap map[string]*FlatStruct
}

type ComplexStruct struct {
	XportVoid   VoidStruct
	XportStruct Struct
	XportPtr    *FlatStruct
	XportSlice  []FlatStruct
	XportArray  [2]FlatStruct
	XportMap    map[string]*FlatStruct
}

var mapTestCases = []MapTestCase{
	{"Void", maker{}, new(VoidStruct), `{}`},
	{"Flat", maker{}, new(FlatStruct), `{"Xport":0}`},
	// Bug for final zero-size field, see https://github.com/golang/go/issues/18016
	// {"AnonymousVoidStruct", maker{}, new(AnonymousVoidStruct), `{"Xport":0}`},
	{"Struct", maker{}, new(Struct), `{"Xport1":0,"Xport2":{"Xport":0}}`},
	{"Ptr", maker{}, new(PtrStruct), `{"Xport1":0,"Xport2":null}`},
	{"Slice", maker{}, new(SliceStruct), `{"XportSlice":null,"XportArray":[{"Xport":0},{"Xport":0}]}`},
	{"Map", maker{}, &MapStruct{XportMap: map[string]*FlatStruct{"A": {Xport: 1}, "B": {Xport: 2}}},
		`{"XportMap":{"A":{"Xport":1},"B":{"Xport":2}}}`},
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
	test.Run("Unsupported", func(test *testing.T) {
		defer shouldPanic(test)
		Convert(new(struct{ I interface{} }), maker{})
	})
}

func shouldPanic(test *testing.T) {
	if p := recover(); p == nil {
		test.Fatal("It should panic")
	}
}

func BenchmarkConvert(b *testing.B) {
	p := new(ComplexStruct)
	b.Run("Cached", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Convert(p, maker{})
		}
	})
	b.Run("Cold", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			cache.m = make(map[cacheKey]reflect.Type)
			b.StartTimer()
			Convert(p, maker{})
		}
	})
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

// func TestX(test *testing.T) {
// 	reflect.StructOf([]reflect.StructField{
// 		{Name: "a", Type: reflect.TypeOf("")},
// 		{Name: "b", Type: reflect.TypeOf("")},
// 	})
// }

func Example_viewOfData() {
	type UserProfile struct {
		Id          int64
		Name        string `view:"*"`
		CardNumber  string `view:"user"`
		SupportNote string `view:"support"`
	}
	profile := &UserProfile{
		Id:          7,
		Name:        "Duke Nukem",
		CardNumber:  "4378 0990 7823 1019",
		SupportNote: "Strange customer",
	}

	userView := Convert(profile, NewView("json", "user"))
	supportView := Convert(profile, NewView("json", "support"))

	// Now profile, userView and supportView point
	// on the same memory but have different types
	// with different tags.

	b, _ := json.MarshalIndent(userView, "", "  ")
	fmt.Println(string(b))
	b, _ = json.MarshalIndent(supportView, "", "  ")
	fmt.Println(string(b))
	// Output:
	// {
	//   "Name": "Duke Nukem",
	//   "CardNumber": "4378 0990 7823 1019"
	// }
	// {
	//   "Name": "Duke Nukem",
	//   "SupportNote": "Strange customer"
	// }
}
