// dynstruct.go -- dynamically create structs from proto schema

package proto

import (
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"log"
	"reflect"
	"strings"
)

//TODO should we use & for pointer type
var typeNames = map[string]reflect.Type{
	"int32":    reflect.TypeOf(int32(0)),
	"int64":    reflect.TypeOf(int64(0)),
	"float32":  reflect.TypeOf(float32(0)),
	"float64":  reflect.TypeOf(float64(0)),
	"bool":     reflect.TypeOf(true),
	"string":   reflect.TypeOf(""),
	"byte":     reflect.TypeOf(byte(0)),
	"*int32":   reflect.TypeOf(int32(0)),
	"*int64":   reflect.TypeOf(int64(0)),
	"*bool":    reflect.TypeOf(true),
	"*string":  reflect.TypeOf(""),
	"*byte":    reflect.TypeOf(byte(0)),
	"*uint32":  reflect.TypeOf(uint32(0)),
	"*uint64":  reflect.TypeOf(uint64(0)),
	"uint32":   reflect.TypeOf(uint32(0)),
	"uint64":   reflect.TypeOf(uint64(0)),
	"*float32": reflect.TypeOf(float32(0)),
	"*float64": reflect.TypeOf(float64(0)),
}

func DynEnumField(gen *Generator, message *Descriptor, field *descriptor.FieldDescriptorProto, structs *map[string]*reflect.Type) (ref reflect.StructField) {
	return ref
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
	fname := "anyProtoFile.proto"
	gen.Request.FileToGenerate = []string{fname}
	desc.Name = &fname
	gen.Request.ProtoFile = []*descriptor.FileDescriptorProto{desc}
	gen.CommandLineParameters("")
	gen.WrapTypes()
	gen.SetPackageNames()
	gen.BuildTypeNameMap()

	gen.GenerateAllFiles()

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
func DynStructTest(desc *descriptor.FileDescriptorProto) []*reflect.Type {
	structs := make([]*reflect.Type, 0, len(desc.MessageType))
	gen := makeSyntheticGenerator(desc)
	for _, message := range gen.allFiles[0].desc {
		fields := make([]reflect.StructField, 0, len(message.Field))
		for _, field := range message.Field {
			fields = append(fields, DynFieldTest(gen, message, field))
		}
		struct_ := reflect.StructOf(fields)
		structs = append(structs, &struct_)
	}
	return structs
}
func DynFieldTest(gen *Generator, message *Descriptor, field *descriptor.FieldDescriptorProto) reflect.StructField {
	var wire string
	var field_type reflect.Type
	if field.Type != nil {
		if *field.Type == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			fileDescriptorProto := gen.ObjectNamed(*field.TypeName).File()
			tempMessage := fileDescriptorProto.GetMessageType()[0]
			for _, tempMessage = range fileDescriptorProto.GetMessageType() {
				if tempMessage.GetName() == *field.TypeName {
					break
				}
			}

			fields := make([]reflect.StructField, 0, len(tempMessage.Field))
			for _, field := range tempMessage.Field {
				fields = append(fields, DynFieldTest(gen, message, field))
			}
			field_type = reflect.PtrTo(reflect.StructOf(fields))
			wire = "bytes"
			tag := gen.goTag(message, field, wire) //which field to use it here
			if field.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
				field_type = reflect.SliceOf(field_type)
			}
			return reflect.StructField{
				Name: CamelCase(*field.Name),
				Type: field_type,
				Tag:  reflect.StructTag("protobuf:" + tag),
			}
		} else if *field.Type == descriptor.FieldDescriptorProto_TYPE_ENUM {
			fileDescriptorProto := gen.ObjectNamed(*field.TypeName).File()
			enum := fileDescriptorProto.GetEnumType()[0]
			for _, enum = range fileDescriptorProto.GetEnumType() {
				if enum.GetName() == *field.TypeName {
					break
				}
			}
			for _, field2 := range fileDescriptorProto.GetMessageType() {
				if field2.GetName() == *field.TypeName {
					break
				}
			}

			_, wire := gen.GoType(message, field)
			//TODO find a way to find the type of enum using typename return from GoType
			field_type = reflect.TypeOf(int32(0))
			tag := gen.goTag(message, field, wire)

			return reflect.StructField{
				Name: CamelCase(*field.Name),
				Type: field_type,
				Tag:  reflect.StructTag("protobuf:" + tag),
			}
		} else {
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

		}
	} else if field.TypeName != nil {
		log.Panicf("Error TypeName is nil")
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
func UpdateTypeNameWithType(fileDescriptor *descriptor.FileDescriptorProto) {
	for _, message := range fileDescriptor.GetMessageType() {
		//fmt.Println("descriptor proto to update index = ", index, message)
		for _, field := range message.GetField() {
			if field.GetType() == descriptor.FieldDescriptorProto_TYPE_DOUBLE {
				//fmt.Println("set type for ", field.GetName(), field.GetTypeName())
				findType(field, fileDescriptor)
				//fmt.Println("field = ", field)
			}
		}
	}
	return
}
func findType(field *descriptor.FieldDescriptorProto, fileDescriptor *descriptor.FileDescriptorProto) descriptor.FieldDescriptorProto_Type {
    for _, msg := range fileDescriptor.GetMessageType() {
        if msg.GetName() == field.GetTypeName() {
            field.Type = newType(descriptor.FieldDescriptorProto_TYPE_MESSAGE)
            return descriptor.FieldDescriptorProto_TYPE_MESSAGE
        }
        //fmt.Println("msgName ", msg.GetName(), field.GetTypeName())
    }
    field.Type = newType(descriptor.FieldDescriptorProto_TYPE_ENUM)
    return descriptor.FieldDescriptorProto_TYPE_ENUM
}
func newType(inputType descriptor.FieldDescriptorProto_Type) (fieldType *descriptor.FieldDescriptorProto_Type) {
    fieldType = new(descriptor.FieldDescriptorProto_Type)
    *fieldType = inputType
    return
}
