// Code generated by protoc-gen-go. DO NOT EDIT.
// source: doer.proto

package doer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Command struct {
	Thing                string   `protobuf:"bytes,1,opt,name=thing" json:"thing,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Command) Reset()         { *m = Command{} }
func (m *Command) String() string { return proto.CompactTextString(m) }
func (*Command) ProtoMessage()    {}
func (*Command) Descriptor() ([]byte, []int) {
	return fileDescriptor_doer_bd32fdb7b2ef4951, []int{0}
}
func (m *Command) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Command.Unmarshal(m, b)
}
func (m *Command) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Command.Marshal(b, m, deterministic)
}
func (dst *Command) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Command.Merge(dst, src)
}
func (m *Command) XXX_Size() int {
	return xxx_messageInfo_Command.Size(m)
}
func (m *Command) XXX_DiscardUnknown() {
	xxx_messageInfo_Command.DiscardUnknown(m)
}

var xxx_messageInfo_Command proto.InternalMessageInfo

func (m *Command) GetThing() string {
	if m != nil {
		return m.Thing
	}
	return ""
}

type Response struct {
	Words                string   `protobuf:"bytes,1,opt,name=words" json:"words,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_doer_bd32fdb7b2ef4951, []int{1}
}
func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (dst *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(dst, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetWords() string {
	if m != nil {
		return m.Words
	}
	return ""
}

func init() {
	proto.RegisterType((*Command)(nil), "doer.Command")
	proto.RegisterType((*Response)(nil), "doer.Response")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// DoerClient is the client API for Doer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DoerClient interface {
	DoIt(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Response, error)
	KeepDoing(ctx context.Context, opts ...grpc.CallOption) (Doer_KeepDoingClient, error)
}

type doerClient struct {
	cc *grpc.ClientConn
}

func NewDoerClient(cc *grpc.ClientConn) DoerClient {
	return &doerClient{cc}
}

func (c *doerClient) DoIt(ctx context.Context, in *Command, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/doer.Doer/DoIt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *doerClient) KeepDoing(ctx context.Context, opts ...grpc.CallOption) (Doer_KeepDoingClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Doer_serviceDesc.Streams[0], "/doer.Doer/KeepDoing", opts...)
	if err != nil {
		return nil, err
	}
	x := &doerKeepDoingClient{stream}
	return x, nil
}

type Doer_KeepDoingClient interface {
	Send(*Command) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type doerKeepDoingClient struct {
	grpc.ClientStream
}

func (x *doerKeepDoingClient) Send(m *Command) error {
	return x.ClientStream.SendMsg(m)
}

func (x *doerKeepDoingClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DoerServer is the server API for Doer service.
type DoerServer interface {
	DoIt(context.Context, *Command) (*Response, error)
	KeepDoing(Doer_KeepDoingServer) error
}

func RegisterDoerServer(s *grpc.Server, srv DoerServer) {
	s.RegisterService(&_Doer_serviceDesc, srv)
}

func _Doer_DoIt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Command)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoerServer).DoIt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/doer.Doer/DoIt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoerServer).DoIt(ctx, req.(*Command))
	}
	return interceptor(ctx, in, info, handler)
}

func _Doer_KeepDoing_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DoerServer).KeepDoing(&doerKeepDoingServer{stream})
}

type Doer_KeepDoingServer interface {
	Send(*Response) error
	Recv() (*Command, error)
	grpc.ServerStream
}

type doerKeepDoingServer struct {
	grpc.ServerStream
}

func (x *doerKeepDoingServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *doerKeepDoingServer) Recv() (*Command, error) {
	m := new(Command)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Doer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "doer.Doer",
	HandlerType: (*DoerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DoIt",
			Handler:    _Doer_DoIt_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "KeepDoing",
			Handler:       _Doer_KeepDoing_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "doer.proto",
}

func init() { proto.RegisterFile("doer.proto", fileDescriptor_doer_bd32fdb7b2ef4951) }

var fileDescriptor_doer_bd32fdb7b2ef4951 = []byte{
	// 145 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0xc9, 0x4f, 0x2d,
	0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0xe4, 0xb9, 0xd8, 0x9d, 0xf3,
	0x73, 0x73, 0x13, 0xf3, 0x52, 0x84, 0x44, 0xb8, 0x58, 0x4b, 0x32, 0x32, 0xf3, 0xd2, 0x25, 0x18,
	0x15, 0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0x25, 0x05, 0x2e, 0x8e, 0xa0, 0xd4, 0xe2, 0x82, 0xfc,
	0xbc, 0xe2, 0x54, 0x90, 0x8a, 0xf2, 0xfc, 0xa2, 0x94, 0x62, 0x98, 0x0a, 0x30, 0xc7, 0x28, 0x91,
	0x8b, 0xc5, 0x25, 0x3f, 0xb5, 0x48, 0x48, 0x1d, 0x44, 0x7b, 0x96, 0x08, 0xf1, 0xea, 0x81, 0x6d,
	0x81, 0x1a, 0x2b, 0xc5, 0x07, 0xe1, 0xc2, 0x0c, 0x51, 0x62, 0x10, 0x32, 0xe0, 0xe2, 0xf4, 0x4e,
	0x4d, 0x2d, 0x70, 0xc9, 0xcf, 0xcc, 0x4b, 0x27, 0xa8, 0x5a, 0x83, 0xd1, 0x80, 0x31, 0x89, 0x0d,
	0xec, 0x64, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x89, 0xd5, 0xa8, 0x89, 0xc0, 0x00, 0x00,
	0x00,
}