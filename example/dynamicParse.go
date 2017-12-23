package main

import (
	"SampleMsgQ"
	"encoding/json"
	"fmt"
	dynamicProto "github.com/devesh2408/dynproto"
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"io/ioutil"
	"reflect"
)

func main() {
	protomessage := getSampleMessage()
	pbBytes, _ := proto.Marshal(protomessage)

	b, err := ioutil.ReadFile("SampleMsgQ.proto") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
    b2, err := ioutil.ReadFile("JitMsgQInclude.proto") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}


	protoString := []string{string(b), string(b2)} // convert content to a 'string'
//	protoString = protoString + string(b2) // convert content to a 'string'
	//fmt.Println(protoString) // print the content as a 'string'
    fmt.Println("array of protostring = \n", protoString)
	fileDescriptor, err := dynamicProto.ParseMultipleSchema(protoString)
	if err != nil {
		fmt.Println("unmarshal error", err)
		return
	}
    //fmt.Println("file descriptor :", fileDescriptor)
    //fmt.Println(" end of file descriptor :")

	//dyn := dynamicProto.DynStruct(&fileDescriptor)
	UpdateTypeNameWithType(&fileDescriptor)
	msgIndex := 0
	for index, message := range fileDescriptor.GetMessageType() {
		if message.GetName() == "SampleMsgQ" {
			fmt.Println("message name , index ", message.GetName(), index)
            msgIndex = index
			break
		}
	}
	dyn := dynamicProto.DynStructTest(&fileDescriptor)
	//todo: find the index for structure to parse
	val := reflect.New(*dyn[msgIndex])
	buf := dynamicProto.NewBuffer(pbBytes)

	err = buf.Unmarshal(val.Interface())
	if err != nil {
		fmt.Println("unmarshal error")
		return
	}
	jsonValue, _ := json.Marshal(val.Interface())
	fmt.Printf("json Value = %s\n", jsonValue)

	fmt.Println("end of test")
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
func getSampleMessage() *SampleMsgQ.SampleMsgQ {
	evt := SampleMsgQ.Event_SERVICE_ADD
	Samplemsg := &SampleMsgQ.SampleMsgQ{
		Version:               proto.Uint32(1),
		Event:                 &evt,
		ServiceType:           proto.String("BASICSVC"),
		Allowance:             proto.Int64(1000),
		Priority:              proto.Uint32(1),
		HostURL:      proto.String("https://www.google.com"),
		Mdn:          proto.Uint64(9876543210),
	}

	return Samplemsg
}
