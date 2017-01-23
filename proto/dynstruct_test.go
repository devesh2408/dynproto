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

func makeDynamicStruct(a int, b string) ([]*reflect.Type, reflect.Value) {
	schema, _ := ParseSchema("syntax = \"proto3\"; message Hello {int32 a = 1; string b = 3;}")
	dyn := DynStruct(&schema)
	val := reflect.New(*dyn[0])
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
	val2 := reflect.New(*dyn[0])
	buf.Unmarshal(val2.Interface())
	if val2.Elem().FieldByName("A").Int() != 25 || val.Elem().FieldByName("B").String() != "twenty-five" {
		t.Error("readback")
	}
}

func TestDynEnum(t *testing.T) { t.Fail() }

func TestDynRepeated(t *testing.T) {
	schema, _ := ParseSchema(`syntax = "proto3"; message Hello {repeated int32 a = 1; string b = 2;}`)
	dyn := DynStruct(&schema)
	val := reflect.New(*dyn[0])
	val.Elem().FieldByName("A").Set(reflect.ValueOf([]int32{1, 2}))
	val.Elem().FieldByName("B").SetString("Hello.b")

	buf := Buffer{}
	buf.Marshal(val.Interface())
	readback := reflect.New(*dyn[0])
	buf.Unmarshal(readback.Interface())
	if readback.Elem().FieldByName("A").Len() != 2 || readback.Elem().FieldByName("A").Index(0).Int() != 1 {
		t.Error("readback slice")
	}
}

func TestDynMap(t *testing.T) { t.Fail() }

func TestDynNested(t *testing.T) {
	schema, _ := ParseSchema(`syntax = "proto3";
		message Nested {int32 a = 1; string b = 2;}
		message Hello {int32 c = 1; Nested d = 2;}
	`)
	dyn := DynStruct(&schema)
	val0 := reflect.New(*dyn[0])
	val0.Elem().FieldByName("A").SetInt(100)
	val0.Elem().FieldByName("B").SetString("nested_b")
	val1 := reflect.New(*dyn[1])
	val1.Elem().FieldByName("C").SetInt(200)
	val1.Elem().FieldByName("D").Set(val0)

	buf := Buffer{}
	buf.Marshal(val1.Interface())
	readback := reflect.New(*dyn[1])
	buf.Unmarshal(readback.Interface())
	inner := readback.Elem().FieldByName("D").Elem()
	if readback.Elem().FieldByName("C").Int() != 200 || inner.FieldByName("A").Int() != 100 || inner.FieldByName("B").String() != "nested_b" {
		t.Error("readback")
	}
}
