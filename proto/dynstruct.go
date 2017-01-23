// dynstruct.go -- dynamically create structs from proto schema

package proto

import (
	"log"
	"reflect"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var typeNames = map[string]reflect.Type{
	"int32":   reflect.TypeOf(int32(0)),
	"int64":   reflect.TypeOf(int64(0)),
	"float32": reflect.TypeOf(float32(0)),
	"float64": reflect.TypeOf(float64(0)),
	"bool":    reflect.TypeOf(true),
	"string":  reflect.TypeOf(""),
	"byte":    reflect.TypeOf(byte(0)),
}

func DynField(gen *Generator, message *Descriptor, field *descriptor.FieldDescriptorProto, structs *map[string]*reflect.Type) reflect.StructField {
	var wire string
	var field_type reflect.Type
	if field.Type != nil {
		var typename string
		typename, wire = gen.GoType(message, field)
		is_slice := strings.HasPrefix(typename, "[]")
		var ok bool
		if is_slice {
			field_type, ok = typeNames[typename[2:]]
		} else {
			field_type, ok = typeNames[typename]
		}
		if !ok {
			log.Panicf("unk typename %s", typename)
		}
		if is_slice {
			field_type = reflect.SliceOf(field_type)
		}
	} else if field.TypeName != nil {
		if nested, ok := (*structs)[*field.TypeName]; ok {
			field_type = reflect.PtrTo(*nested)
			wire = "bytes"
			// note: not sure why this isn't being set
			type_ := descriptor.FieldDescriptorProto_TYPE_MESSAGE
			field.Type = &type_
		} else {
			log.Panicf("embedded field %v not in structs", *field.TypeName)
		}
	} else {
		panic("need Type or TypeName")
	}
	tag := gen.goTag(message, field, wire)
	return reflect.StructField{
		Name: CamelCase(*field.Name),
		Type: field_type,
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
func DynStruct(desc *descriptor.FileDescriptorProto) []*reflect.Type {
	structs := make([]*reflect.Type, 0, len(desc.MessageType))
	structByName := make(map[string]*reflect.Type)
	gen := makeSyntheticGenerator(desc)
	for _, message := range gen.allFiles[0].desc {
		fields := make([]reflect.StructField, 0, len(message.Field))
		for _, field := range message.Field {
			fields = append(fields, DynField(gen, message, field, &structByName))
		}
		struct_ := reflect.StructOf(fields)
		structs = append(structs, &struct_)
		structByName[*message.Name] = &struct_
	}
	return structs
}
