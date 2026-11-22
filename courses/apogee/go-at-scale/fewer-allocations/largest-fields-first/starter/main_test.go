package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

// fields are glance's fields and their types, in the order the strip first had them.
var fields = []struct{ name, typ string }{
	{"Sat", "uint8"},
	{"Station", "uint8"},
	{"AOS", "int64"},
	{"LOS", "int64"},
	{"MaxElevDeg", "float32"},
	{"Dish", "uint8"},
	{"Booked", "bool"},
}

// layout says where each of glance's fields starts, in the order they're declared.
func layout() string {
	t := reflect.TypeFor[glance]()
	var parts []string
	for f := range t.Fields() {
		parts = append(parts, fmt.Sprintf("%s at %d", f.Name, f.Offset))
	}
	return strings.Join(parts, ", ")
}

func TestSameFields(t *testing.T) {
	ty := reflect.TypeFor[glance]()
	for _, want := range fields {
		f, ok := ty.FieldByName(want.name)
		switch {
		case !ok:
			t.Errorf("glance has no field %s, want %s %s", want.name, want.name, want.typ)
		case f.Type.String() != want.typ:
			t.Errorf("glance.%s is %s, want %s", want.name, f.Type, want.typ)
		}
	}
	if ty.NumField() != len(fields) {
		t.Errorf("glance has %d fields, want the strip's %d", ty.NumField(), len(fields))
	}
}

func TestNoPadding(t *testing.T) {
	if got := unsafe.Sizeof(glance{}); got != 24 {
		t.Errorf("unsafe.Sizeof(glance{}) = %d, want 24, its fields' sum (%s)", got, layout())
	}
}
