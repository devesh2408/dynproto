// cwrapper.cpp -- C wrapper for C++ proto APIs
// todo: look into protobuf-c as a way to get this for free

#include "cwrapper.h"
#include <google/protobuf/io/zero_copy_stream_impl_lite.h>
#include <google/protobuf/io/tokenizer.h>
#include <google/protobuf/compiler/parser.h>

extern "C" {

struct FakeString parse_schema(const char* data, int size){
    using namespace google::protobuf;
    io::ArrayInputStream stream((const void*)data, size);
    // ErrorCollector ec; // todo: surface parse errors to go
    io::Tokenizer tokenizer(&stream, NULL);
    compiler::Parser parser;
    FileDescriptorProto fdp;
    FakeString fs = {NULL, 0};
    if (!parser.Parse(&tokenizer, &fdp)) return fs;
    std::string blob;
    if (!fdp.SerializeToString(&blob)) return fs;
    void* buf = malloc(blob.size());
    memcpy(buf, blob.data(), blob.size());
    fs.bytes = buf;
    fs.size = blob.size();
    return fs;
}

} // extern "C"
