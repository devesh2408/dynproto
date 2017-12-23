// wrapper.go -- wrapper around libprotobuf

package proto

import (
	"errors"
    "unsafe"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// #cgo pkg-config: protobuf
// #include <stdlib.h>
// #include "cwrapper.h"
import "C"

var SchemaParseError = errors.New("error parsing schema")

// wrapper for CPP proto schema parser
func ParseSchema(schema string) (descriptor.FileDescriptorProto, error) {
	desc := descriptor.FileDescriptorProto{}
	blob := C.parse_schema(C.CString(schema), C.int(len(schema)))
	if uintptr(blob.bytes) == 0 {
		return desc, SchemaParseError
	}
	defer C.free(blob.bytes)
	raw := C.GoBytes(blob.bytes, C.int(blob.size))
	if err := proto.Unmarshal(raw, &desc); err != nil {
		return desc, err
	} else {
		return desc, nil
	}
}
func ParseMultipleSchema(schema []string) (descriptor.FileDescriptorProto, error) {
    protobufMessages := (*C.ProtobufMessage)(C.malloc( C.sizeof_ProtobufMessage));
    defer C.free(unsafe.Pointer(protobufMessages))
    schemaLen := 0
    for i:=0; i< len(schema); i++ {
        protobufMessages.protoFile[i] = C.CString(schema[i])
        defer C.free(unsafe.Pointer(protobufMessages.protoFile[i]));
        protobufMessages.size[i] = C.int(len(schema[i]))
        schemaLen  += 1

    }

	desc := descriptor.FileDescriptorProto{}
	//blob := C.parse_multiple_schema(C.CString(schema[0]), C.int(len(schema[0])), C.CString(schema[1]), C.int(len(schema[1])))
    blob := C.parseArraySchema((*C.ProtobufMessage)(unsafe.Pointer(protobufMessages)), C.int(schemaLen))
	if uintptr(blob.bytes) == 0 {
		return desc, SchemaParseError
	}
	defer C.free(blob.bytes)
	raw := C.GoBytes(blob.bytes, C.int(blob.size))
	if err := proto.Unmarshal(raw, &desc); err != nil {
		return desc, err
	} else {
		return desc, nil
	}
}

