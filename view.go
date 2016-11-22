package retag

import (
	"fmt"
	"reflect"
	"strings"
)

func NewView(tag, name string) TagMaker {
	return tagView{name, tag}
}

type tagView struct {
	name string
	tag  string
}

func (v tagView) MakeTag(t reflect.Type, fieldIndex int) reflect.StructTag {
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

func (v tagView) isMatch(tag string) bool {
	if tag == "*" {
		return true
	}
	list := parseStringList(tag)
	if list.contains(v.name) {
		return true
	}
	return false
}

func parseStringList(list string) stringList {
	return strings.Split(list, ",")
}

type stringList []string

func (l stringList) contains(s string) bool {
	for _, i := range l {
		if i == s {
			return true
		}
	}
	return false
}
