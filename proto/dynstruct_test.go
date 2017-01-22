package proto

import (
	"reflect"
	"testing"
)

func TestParseSchema(t *testing.T) {
	schema, err := ParseSchema("syntax = \"proto3\"; message Hello {int32 x = 1;}")
	if err != nil {
		t.Error(err)
	}
	if len(schema.MessageType) != 1 {
		t.Error("len(types)")
	}
	desc := schema.MessageType[0]
	if *desc.Name != "Hello" {
		t.Error("name")
	}
	// test unparseable
	_, err = ParseSchema("syntax = \"proto2\"; message NoSuchType {int x = 1;}")
	if err == nil {
		t.Error("unexpected success")
	}
}

func makeDynamicStruct(a int, b string) ([]reflect.Type, reflect.Value) {
	schema, _ := ParseSchema("syntax = \"proto3\"; message Hello {int32 a = 1; string b = 3;}")
	dyn := DynStruct(&schema)
	val := reflect.New(dyn[0])
	val.Elem().FieldByName("A").SetInt(int64(a))
	val.Elem().FieldByName("B").SetString(b)
	return dyn, val
}

type EmptyStruct struct {
	a int
}
type fntype func(int) error

func TestDynStruct(t *testing.T) {
	_, val := makeDynamicStruct(25, "twenty-five")
	if val.Elem().FieldByName("A").Int() != 25 || val.Elem().FieldByName("B").String() != "twenty-five" {
		t.Error("readback")
	}
}

func TestDynMarshalUnmarshal(t *testing.T) {
	dyn, val := makeDynamicStruct(25, "twenty-five")
	buf := Buffer{}
	buf.Marshal(val.Interface())
	val2 := reflect.New(dyn[0])
	buf.Unmarshal(val2.Interface())
	if val2.Elem().FieldByName("A").Int() != 25 || val.Elem().FieldByName("B").String() != "twenty-five" {
		t.Error("readback")
	}
}
