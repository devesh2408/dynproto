## dynproto

This is a fork of github.com/golang/protobuf to add schema parsing & marshaling/unmarshaling of blobs from dynamically parsed schemas.

Use this to parse protos when schemas aren't available at compile-time.

### Example:

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

### Changes from upstream:
* use cgo to wrap schema parsing from libprotobuf.
* added DynProto function which returns a reflect.Type (a dynamic struct)
* Buffer.Marshal / Unmarshal takes interface{} instead of proto.Message
* proto package copies parts of `protoc-gen-go/generator.go`.

### Status & concerns:
* apparently working: proto messages with scalars and strings
* likely not working: enums and nested protos. Also not sure about repeated fields & maps. (if you want to fix, look at the `panic`s in `proto/generator.go`).
* I don't know how to test this project under all build constraints (appengine / JS). I assume it fails for the JS target because of cgo.
* cgo piece (schemaparser.go) needs review from someone who understands cgo to check for memory leaks
* haven't done any performance testing. Should be way slower than generated marshal/unmarshal functions.
* security: I don't know if the schema parsing & ser/des code were written with user-provided input in mind. You should pay for a security audit before pointing this at the internet.

### Future work:
* I'm evaluating this for use in a larger project. If that moves forward, I plan to clean up this fork and ask the maintainers to merge back my changes.
* Don't expect this repo to stay functional or up-to-date. If your organization wants to use these features, you should (a) write more tests (b) benchmark and (c) take over the project.
