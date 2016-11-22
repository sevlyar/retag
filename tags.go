package retag

import (
	// "fmt"
	"reflect"
	"sync"
	"unsafe"
)

type TagMaker interface {
	MakeTag(t reflect.Type, fieldIndex int) reflect.StructTag
}

// Doesn't support cyclic dependencies
func Convert(p interface{}, maker TagMaker) interface{} {
	strPtrVal := reflect.ValueOf(p)
	// TODO: check type
	newType := getType(strPtrVal.Type().Elem(), maker)
	newPtrVal := reflect.NewAt(newType, unsafe.Pointer(strPtrVal.Pointer()))
	return newPtrVal.Interface()
}

type cacheKey struct {
	reflect.Type
	TagMaker
}

var cache = struct {
	sync.RWMutex
	m map[cacheKey]reflect.Type
}{
	m: make(map[cacheKey]reflect.Type),
}

func getType(structType reflect.Type, maker TagMaker) reflect.Type {
	key := cacheKey{structType, maker}
	cache.RLock()
	t, ok := cache.m[key]
	cache.RUnlock()
	if !ok {
		t = makeType(structType, maker)
		cache.Lock()
		cache.m[key] = t
		cache.Unlock()
	}
	return t
}

// TODO: optimize cases when type is not modified
func makeType(t reflect.Type, maker TagMaker) reflect.Type {
	switch t.Kind() {
	case reflect.Struct:
		return makeStructType(t, maker)
	case reflect.Ptr:
		return reflect.PtrTo(getType(t.Elem(), maker))
	case reflect.Array:
		return reflect.ArrayOf(t.Len(), getType(t.Elem(), maker))
	case reflect.Slice:
		return reflect.SliceOf(getType(t.Elem(), maker))
	case reflect.Map:
		return reflect.MapOf(getType(t.Key(), maker), getType(t.Elem(), maker))
	case
		reflect.Chan,
		reflect.Func,
		reflect.UnsafePointer,
		reflect.Interface:
		panic("tags.Map: Unsupported type: " + t.Kind().String())
	default:
		// don't modify type in another case
		return t
	}
}

func makeStructType(structType reflect.Type, maker TagMaker) reflect.Type {
	if structType.NumField() == 0 {
		return structType
	}

	fields := make([]reflect.StructField, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		strField := structType.Field(i)
		if isExported(strField.Name) {
			strField.Type = getType(strField.Type, maker)
		} else {
			// strange case with anonymous fields
			strField.PkgPath = ""
			strField.Name = ""
			// strField.Anonymous = true
		}
		strField.Tag = maker.MakeTag(structType, i)
		fields = append(fields, strField)
	}

	newType := reflect.StructOf(fields)
	if structType.Size() != newType.Size() {
		// TODO: debug
		// fmt.Println(newType.Size(), newType)
		// for i := 0; i < newType.NumField(); i++ {
		// 	fmt.Println(newType.Field(i))
		// }
		// fmt.Println(structType.Size(), structType)
		// for i := 0; i < structType.NumField(); i++ {
		// 	fmt.Println(structType.Field(i))
		// }
		panic("tags.Map: Unexpected case - type has a size different from size of original type")
	}
	return newType
}

func isExported(name string) bool {
	b := name[0]
	return !('a' <= b && b <= 'z') && b != '_'
}
