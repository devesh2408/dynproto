// wrapper.go -- wrapper around libprotobuf

package proto

import (
	"errors"

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
