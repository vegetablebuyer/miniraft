// Code generated by protoc-gen-go.
// source: logentry.proto
// DO NOT EDIT!

/*
Package protobuf is a generated protocol buffer package.

It is generated from these files:
	logentry.proto

It has these top-level messages:
	LogEntry
*/
package protobuf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type LogEntry struct {
	Index       uint64 `protobuf:"varint,1,opt,name=Index" json:"Index,omitempty"`
	Term        uint64 `protobuf:"varint,2,opt,name=Term" json:"Term,omitempty"`
	CommandName string `protobuf:"bytes,3,opt,name=CommandName" json:"CommandName,omitempty"`
	Command     []byte `protobuf:"bytes,4,opt,name=Command,proto3" json:"Command,omitempty"`
}

func (m *LogEntry) Reset()                    { *m = LogEntry{} }
func (m *LogEntry) String() string            { return proto.CompactTextString(m) }
func (*LogEntry) ProtoMessage()               {}
func (*LogEntry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*LogEntry)(nil), "protobuf.LogEntry")
}

func init() { proto.RegisterFile("logentry.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 132 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0xc9, 0x4f, 0x4f,
	0xcd, 0x2b, 0x29, 0xaa, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x00, 0x53, 0x49, 0xa5,
	0x69, 0x4a, 0x05, 0x5c, 0x1c, 0x3e, 0xf9, 0xe9, 0xae, 0x20, 0x39, 0x21, 0x11, 0x2e, 0x56, 0xcf,
	0xbc, 0x94, 0xd4, 0x0a, 0x09, 0x46, 0x05, 0x46, 0x0d, 0x96, 0x20, 0x08, 0x47, 0x48, 0x88, 0x8b,
	0x25, 0x24, 0xb5, 0x28, 0x57, 0x82, 0x09, 0x2c, 0x08, 0x66, 0x0b, 0x29, 0x70, 0x71, 0x3b, 0xe7,
	0xe7, 0xe6, 0x26, 0xe6, 0xa5, 0xf8, 0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2b, 0x30, 0x6a, 0x70, 0x06,
	0x21, 0x0b, 0x09, 0x49, 0x70, 0xb1, 0x43, 0xb9, 0x12, 0x2c, 0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x30,
	0x6e, 0x12, 0x1b, 0xd8, 0x6e, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5f, 0xbe, 0x7b, 0x3f,
	0x94, 0x00, 0x00, 0x00,
}
