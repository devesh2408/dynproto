// dynstruct.go -- dynamically create structs from proto schema

package proto

import (
	"reflect"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func DynType(desc *descriptor.FieldDescriptorProto) reflect.Type {
	switch *desc.Type.Enum() {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		return reflect.TypeOf(float64(0))
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return reflect.TypeOf(float32(0))
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		return reflect.TypeOf(int64(0))
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		return reflect.TypeOf(uint64(0))
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		return reflect.TypeOf(int32(0))
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		panic("fixed")
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return reflect.TypeOf(true)
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return reflect.TypeOf("")
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		panic("group")
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		panic("submessage")
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return reflect.TypeOf([]byte(""))
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		return reflect.TypeOf(uint32(0))
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		// dec, cast = "b.DecodeVarint()", fieldTypes[field]
		panic("enum")
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		panic("fixed")
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		return reflect.TypeOf(int32(0))
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		return reflect.TypeOf(int64(0))
	default:
		panic("unk type")
	}
}

func DynField(gen *Generator, message *Descriptor, field *descriptor.FieldDescriptorProto) reflect.StructField {
	typename, wire := gen.GoType(message, field)
	tag := gen.goTag(message, field, wire)
	return reflect.StructField{
		Name: CamelCase(*field.Name),
		Type: DynType(field),
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
