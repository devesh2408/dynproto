// cwrapper.cpp -- C wrapper for C++ proto APIs
// todo: look into protobuf-c as a way to get this for free

#include "cwrapper.h"
#include <google/protobuf/io/zero_copy_stream_impl_lite.h>
#include <google/protobuf/io/tokenizer.h>
#include <google/protobuf/compiler/parser.h>

extern "C" {
    using namespace google::protobuf;
    struct FakeString parse_schema(const char* data, int size){
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
    FileDescriptorProto getFdpFromSchema(const char* data, int size){
        using namespace google::protobuf;
        io::ArrayInputStream stream((const void*)data, size);
        // ErrorCollector ec; // todo: surface parse errors to go
        io::Tokenizer tokenizer(&stream, NULL);
        compiler::Parser parser;
        FileDescriptorProto fdp;
        if (!parser.Parse(&tokenizer, &fdp)) return fdp; //do error return 
        return fdp;
    }


    struct FakeString parse_multiple_schema(const char* data, int size, const char* data2, int size2){
        using namespace google::protobuf;
        io::ArrayInputStream stream((const void*)data, size);
        //printf("sizes = %d %d\n strlen = %ld %ld",size, size2, strlen(data), strlen(data2)); 
        // ErrorCollector ec; // todo: surface parse errors to go
        io::Tokenizer tokenizer(&stream, NULL);
        compiler::Parser parser;
        FileDescriptorProto fdp;
        FakeString fs = {NULL, 0};
        if (!parser.Parse(&tokenizer, &fdp)) return fs;

        io::ArrayInputStream stream2((const void*)data2, size2);
        // ErrorCollector ec; // todo: surface parse errors to go
        io::Tokenizer tokenizer2(&stream2, NULL);
        compiler::Parser parser2;
        FileDescriptorProto fdp2;
        FakeString fs2 = {NULL, 0};
        if (!parser.Parse(&tokenizer2, &fdp2)) return fs2;
        fdp.MergeFrom(fdp2);



        std::string blob;
        if (!fdp.SerializeToString(&blob)) return fs;
        void* buf = malloc(blob.size());
        memcpy(buf, blob.data(), blob.size());
        fs.bytes = buf;
        fs.size = blob.size();
        return fs;

    }
    
    struct FakeString parseArraySchema(ProtobufMessage *in, int arrayLen){
        FileDescriptorProto fdp;
        fdp = getFdpFromSchema(in->protoFile[0], in->size[0]);   
        FakeString fs = {NULL, 0};

        if(arrayLen == 1){
            std::string blob;
            if (!fdp.SerializeToString(&blob)) return fs;
            void* buf = malloc(blob.size());
            memcpy(buf, blob.data(), blob.size());
            fs.bytes = buf;
            fs.size = blob.size();
            return fs;

        }
        for(int i = 1; i < arrayLen; i++)  
        {
            FileDescriptorProto fdp2 = getFdpFromSchema(in->protoFile[i], in->size[i]);   
            fdp.MergeFrom(fdp2);
        }
        std::string blob;
        if (!fdp.SerializeToString(&blob)) return fs;
        void* buf = malloc(blob.size());
        memcpy(buf, blob.data(), blob.size());
        fs.bytes = buf;
        fs.size = blob.size();
        return fs;
    }

} // extern "C"
