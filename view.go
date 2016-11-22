package retag

import (
	"fmt"
	"reflect"
	"strings"
)

func NewView(tag, name string) View {
	return View{name, tag}
}

type View struct {
	name string
	tag  string
}

func (v View) MakeTag(t reflect.Type, fieldIndex int) reflect.StructTag {
	field := t.Field(fieldIndex)
	if v.isMatch(field.Tag.Get("view")) {
		defaultValue := field.Tag.Get(v.tag)
		if defaultValue != "" {
			defaultValue = fmt.Sprintf(`%s:"%s"`, v.tag, defaultValue)
		}
		return reflect.StructTag(defaultValue)
	}
	return reflect.StructTag(v.tag + `:"-"`)
}

func (v View) isMatch(tag string) bool {
	if tag == "*" {
		return true
	}
	list := ParseStringList(tag)
	if list.Contains(v.name) {
		return true
	}
	return false
}

func ParseStringList(list string) StringList {
	return strings.Split(list, ",")
}

type StringList []string

func (l StringList) Contains(s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}
