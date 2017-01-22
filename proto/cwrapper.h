#ifdef __cplusplus
extern "C" {
#endif

struct FakeString {
    void* bytes;
    long unsigned int size;
};

// this returns a serialized blob of the FileDescriptorProto
struct FakeString parse_schema(const char* data, int size);

#ifdef __cplusplus
} // extern "C"
#endif
