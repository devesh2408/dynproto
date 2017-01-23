// dynstruct.go -- dynamically create structs from proto schema

package proto

import (
	"reflect"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var typeNames = map[string]reflect.Type{
	"int32":   reflect.TypeOf(int32(0)),
	"int64":   reflect.TypeOf(int64(0)),
	"float32": reflect.TypeOf(float32(0)),
	"float64": reflect.TypeOf(float64(0)),
	"bool":    reflect.TypeOf(true),
	"string":  reflect.TypeOf(""),
	"[]byte":  reflect.TypeOf([]byte("")),
}

func DynField(gen *Generator, message *Descriptor, field *descriptor.FieldDescriptorProto) reflect.StructField {
	typename, wire := gen.GoType(message, field)
	tag := gen.goTag(message, field, wire)
	return reflect.StructField{
		Name: CamelCase(*field.Name),
		Type: typeNames[typename],
		Tag:  reflect.StructTag("protobuf:" + tag),
	}
}

// this runs the top half of main.go from the generator package so we can access gen.GoType and gen.goTag for fields
func makeSyntheticGenerator(desc *descriptor.FileDescriptorProto) *Generator {
	gen := NewGenerator()
	fname := "synthetic"
	gen.Request.FileToGenerate = []string{fname}
	desc.Name = &fname
	gen.Request.ProtoFile = []*descriptor.FileDescriptorProto{desc}
	gen.WrapTypes()
	gen.SetPackageNames()
	gen.BuildTypeNameMap()
	return gen
}

// warning: we set FileDescriptorProto.Name
func DynStruct(desc *descriptor.FileDescriptorProto) []reflect.Type {
	structs := make([]reflect.Type, 0, len(desc.MessageType))
	gen := makeSyntheticGenerator(desc)
	for _, message := range gen.allFiles[0].desc {
		fields := make([]reflect.StructField, 0, len(message.Field))
		for _, field := range message.Field {
			fields = append(fields, DynField(gen, message, field))
		}
		structs = append(structs, reflect.StructOf(fields))
	}
	return structs
}
