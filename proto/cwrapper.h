#ifdef __cplusplus
extern "C" {
#endif
#define MAX_INCLUDE_PROTO_FILE 5
struct FakeString {
    void* bytes;
    long unsigned int size;
};
//max number of include file is hardcoded to 5 here
typedef struct {
        char *protoFile[MAX_INCLUDE_PROTO_FILE];
        int size[MAX_INCLUDE_PROTO_FILE];
} ProtobufMessage;
//struct FileDescriptorProto;
// this returns a serialized blob of the FileDescriptorProto
struct FakeString parse_schema(const char* data, int size);
struct FakeString parse_multiple_schema(const char* data, int size, const char* data2, int size2);
//struct FileDescriptorProto getFdpFromSchema(const char* data, int size);

struct FakeString parseArraySchema(ProtobufMessage *in, int len);
#ifdef __cplusplus
} // extern "C"
#endif
