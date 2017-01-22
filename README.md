This is a fork of github.com/golang/protobuf to add schema parsing & marshaling/unmarshaling of blobs from dynamically parsed schemas.

Example:

```golang
// careful: this is paraphrased from test suite. Never tested as-is.

schema, _ := proto.ParseSchema("syntax = \"proto3\"; message Hello {int32 a = 1; string b = 3;}")
dyn := proto.DynStruct(&schema)

val := reflect.New(dyn[0]) // dyn[0] i.e. the 0th message in the schema
val.Elem().FieldByName("A").SetInt(10)
val.Elem().FieldByName("B").SetString("ten")

buf := proto.Buffer{}
buf.Marshal(val.Interface())
val2 := reflect.New(dyn[0])
buf.Unmarshal(val2.Interface())

if val2.Elem().FieldByName("A").Int() != 10 || val.Elem().FieldByName("B").String() != "ten" {
	log.Panicf("%v\n", val2)
}
```

Changes from upstream:
* use cgo to wrap schema parsing from libprotobuf.
* Buffer.Marshal / Unmarshal takes interface{} instead of proto.Message
* DynProto function returns a reflect.Type
* proto package copies parts of `protoc-gen-go/generator.go`.

Status & concerns:
* likely not working: enums and nested protos
* I don't know how to test this project under all build constraints (appengine / JS). I assume it fails for the JS target because of cgo.
* cgo piece (schemaparser.go) needs review from someone who understands cgo to check for memory leaks
* haven't done any performance testing. Should be way slower than generated marshal/unmarshal functions.
