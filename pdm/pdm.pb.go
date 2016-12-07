// Code generated by protoc-gen-go.
// source: pdm.proto
// DO NOT EDIT!

/*
Package pdm is a generated protocol buffer package.

It is generated from these files:
	pdm.proto

It has these top-level messages:
	Endpoint
	Handle
	ActionItem
	ActionStatus
	Empty
*/
package pdm

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

type Command int32

const (
	Command_NONE    Command = 0
	Command_ARCHIVE Command = 1
	Command_RESTORE Command = 2
	Command_REMOVE  Command = 3
	Command_CANCEL  Command = 4
)

var Command_name = map[int32]string{
	0: "NONE",
	1: "ARCHIVE",
	2: "RESTORE",
	3: "REMOVE",
	4: "CANCEL",
}
var Command_value = map[string]int32{
	"NONE":    0,
	"ARCHIVE": 1,
	"RESTORE": 2,
	"REMOVE":  3,
	"CANCEL":  4,
}

func (x Command) String() string {
	return proto.EnumName(Command_name, int32(x))
}
func (Command) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Endpoint struct {
	FsUrl   string `protobuf:"bytes,2,opt,name=fs_url,json=fsUrl" json:"fs_url,omitempty"`
	Archive uint32 `protobuf:"varint,1,opt,name=archive" json:"archive,omitempty"`
}

func (m *Endpoint) Reset()                    { *m = Endpoint{} }
func (m *Endpoint) String() string            { return proto.CompactTextString(m) }
func (*Endpoint) ProtoMessage()               {}
func (*Endpoint) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Handle struct {
	Id uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
}

func (m *Handle) Reset()                    { *m = Handle{} }
func (m *Handle) String() string            { return proto.CompactTextString(m) }
func (*Handle) ProtoMessage()               {}
func (*Handle) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type ActionItem struct {
	Id          uint64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Op          Command `protobuf:"varint,2,opt,name=op,enum=pdm.Command" json:"op,omitempty"`
	PrimaryPath string  `protobuf:"bytes,3,opt,name=primary_path,json=primaryPath" json:"primary_path,omitempty"`
	WritePath   string  `protobuf:"bytes,4,opt,name=write_path,json=writePath" json:"write_path,omitempty"`
	Offset      int64   `protobuf:"varint,5,opt,name=offset" json:"offset,omitempty"`
	Length      int64   `protobuf:"varint,6,opt,name=length" json:"length,omitempty"`
	FileId      []byte  `protobuf:"bytes,7,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	Data        []byte  `protobuf:"bytes,8,opt,name=data,proto3" json:"data,omitempty"`
	Uuid        string  `protobuf:"bytes,9,opt,name=uuid" json:"uuid,omitempty"`
	Hash        string  `protobuf:"bytes,10,opt,name=hash" json:"hash,omitempty"`
	Url         string  `protobuf:"bytes,12,opt,name=url" json:"url,omitempty"`
}

func (m *ActionItem) Reset()                    { *m = ActionItem{} }
func (m *ActionItem) String() string            { return proto.CompactTextString(m) }
func (*ActionItem) ProtoMessage()               {}
func (*ActionItem) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type ActionStatus struct {
	Id        uint64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Completed bool    `protobuf:"varint,2,opt,name=completed" json:"completed,omitempty"`
	Error     int32   `protobuf:"varint,3,opt,name=error" json:"error,omitempty"`
	Offset    int64   `protobuf:"varint,4,opt,name=offset" json:"offset,omitempty"`
	Length    int64   `protobuf:"varint,5,opt,name=length" json:"length,omitempty"`
	Handle    *Handle `protobuf:"bytes,6,opt,name=handle" json:"handle,omitempty"`
	FileId    []byte  `protobuf:"bytes,7,opt,name=file_id,json=fileId,proto3" json:"file_id,omitempty"`
	Flags     int32   `protobuf:"varint,8,opt,name=flags" json:"flags,omitempty"`
	Uuid      string  `protobuf:"bytes,9,opt,name=uuid" json:"uuid,omitempty"`
	Hash      string  `protobuf:"bytes,10,opt,name=hash" json:"hash,omitempty"`
	Url       string  `protobuf:"bytes,11,opt,name=url" json:"url,omitempty"`
}

func (m *ActionStatus) Reset()                    { *m = ActionStatus{} }
func (m *ActionStatus) String() string            { return proto.CompactTextString(m) }
func (*ActionStatus) ProtoMessage()               {}
func (*ActionStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *ActionStatus) GetHandle() *Handle {
	if m != nil {
		return m.Handle
	}
	return nil
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func init() {
	proto.RegisterType((*Endpoint)(nil), "pdm.Endpoint")
	proto.RegisterType((*Handle)(nil), "pdm.Handle")
	proto.RegisterType((*ActionItem)(nil), "pdm.ActionItem")
	proto.RegisterType((*ActionStatus)(nil), "pdm.ActionStatus")
	proto.RegisterType((*Empty)(nil), "pdm.Empty")
	proto.RegisterEnum("pdm.Command", Command_name, Command_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for DataMover service

type DataMoverClient interface {
	Register(ctx context.Context, in *Endpoint, opts ...grpc.CallOption) (*Handle, error)
	GetActions(ctx context.Context, in *Handle, opts ...grpc.CallOption) (DataMover_GetActionsClient, error)
	StatusStream(ctx context.Context, opts ...grpc.CallOption) (DataMover_StatusStreamClient, error)
}

type dataMoverClient struct {
	cc *grpc.ClientConn
}

func NewDataMoverClient(cc *grpc.ClientConn) DataMoverClient {
	return &dataMoverClient{cc}
}

func (c *dataMoverClient) Register(ctx context.Context, in *Endpoint, opts ...grpc.CallOption) (*Handle, error) {
	out := new(Handle)
	err := grpc.Invoke(ctx, "/pdm.DataMover/Register", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataMoverClient) GetActions(ctx context.Context, in *Handle, opts ...grpc.CallOption) (DataMover_GetActionsClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DataMover_serviceDesc.Streams[0], c.cc, "/pdm.DataMover/GetActions", opts...)
	if err != nil {
		return nil, err
	}
	x := &dataMoverGetActionsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DataMover_GetActionsClient interface {
	Recv() (*ActionItem, error)
	grpc.ClientStream
}

type dataMoverGetActionsClient struct {
	grpc.ClientStream
}

func (x *dataMoverGetActionsClient) Recv() (*ActionItem, error) {
	m := new(ActionItem)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dataMoverClient) StatusStream(ctx context.Context, opts ...grpc.CallOption) (DataMover_StatusStreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DataMover_serviceDesc.Streams[1], c.cc, "/pdm.DataMover/StatusStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &dataMoverStatusStreamClient{stream}
	return x, nil
}

type DataMover_StatusStreamClient interface {
	Send(*ActionStatus) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type dataMoverStatusStreamClient struct {
	grpc.ClientStream
}

func (x *dataMoverStatusStreamClient) Send(m *ActionStatus) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dataMoverStatusStreamClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for DataMover service

type DataMoverServer interface {
	Register(context.Context, *Endpoint) (*Handle, error)
	GetActions(*Handle, DataMover_GetActionsServer) error
	StatusStream(DataMover_StatusStreamServer) error
}

func RegisterDataMoverServer(s *grpc.Server, srv DataMoverServer) {
	s.RegisterService(&_DataMover_serviceDesc, srv)
}

func _DataMover_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Endpoint)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataMoverServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pdm.DataMover/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataMoverServer).Register(ctx, req.(*Endpoint))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataMover_GetActions_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Handle)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DataMoverServer).GetActions(m, &dataMoverGetActionsServer{stream})
}

type DataMover_GetActionsServer interface {
	Send(*ActionItem) error
	grpc.ServerStream
}

type dataMoverGetActionsServer struct {
	grpc.ServerStream
}

func (x *dataMoverGetActionsServer) Send(m *ActionItem) error {
	return x.ServerStream.SendMsg(m)
}

func _DataMover_StatusStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DataMoverServer).StatusStream(&dataMoverStatusStreamServer{stream})
}

type DataMover_StatusStreamServer interface {
	SendAndClose(*Empty) error
	Recv() (*ActionStatus, error)
	grpc.ServerStream
}

type dataMoverStatusStreamServer struct {
	grpc.ServerStream
}

func (x *dataMoverStatusStreamServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dataMoverStatusStreamServer) Recv() (*ActionStatus, error) {
	m := new(ActionStatus)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _DataMover_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pdm.DataMover",
	HandlerType: (*DataMoverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _DataMover_Register_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetActions",
			Handler:       _DataMover_GetActions_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StatusStream",
			Handler:       _DataMover_StatusStream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("pdm.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 497 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x53, 0xdd, 0x8b, 0xd3, 0x40,
	0x10, 0x37, 0x69, 0x3e, 0x9a, 0x69, 0xee, 0xac, 0x83, 0x1f, 0xa1, 0x28, 0x68, 0x05, 0x39, 0x44,
	0x4e, 0xa9, 0x8f, 0x3e, 0x95, 0x1a, 0xbc, 0x82, 0xd7, 0xca, 0x56, 0x7d, 0x2d, 0xb1, 0xd9, 0xb4,
	0x81, 0xa4, 0x1b, 0x36, 0xdb, 0x93, 0xfb, 0x2f, 0x44, 0xff, 0x61, 0x77, 0x67, 0x7b, 0xb4, 0x8a,
	0xf7, 0xe0, 0xdb, 0xef, 0x63, 0xda, 0x9d, 0xf9, 0xcd, 0x04, 0xa2, 0x26, 0xaf, 0xcf, 0x1b, 0x29,
	0x94, 0xc0, 0x8e, 0x86, 0xc3, 0x77, 0xd0, 0x4d, 0xb7, 0x79, 0x23, 0xca, 0xad, 0xc2, 0x07, 0x10,
	0x14, 0xed, 0x72, 0x27, 0xab, 0xc4, 0x7d, 0xea, 0x9c, 0x45, 0xcc, 0x2f, 0xda, 0x2f, 0xb2, 0xc2,
	0x04, 0xc2, 0x4c, 0xae, 0x36, 0xe5, 0x15, 0x4f, 0x1c, 0xad, 0x9f, 0xb0, 0x1b, 0x3a, 0x4c, 0x20,
	0xb8, 0xc8, 0xb6, 0x79, 0xc5, 0xf1, 0x14, 0xdc, 0x32, 0x27, 0xdb, 0x63, 0x1a, 0x0d, 0x7f, 0xb9,
	0x00, 0xe3, 0x95, 0x2a, 0xc5, 0x76, 0xaa, 0x78, 0xfd, 0xb7, 0x8d, 0x8f, 0xc1, 0x15, 0x0d, 0xbd,
	0x72, 0x3a, 0x8a, 0xcf, 0x4d, 0x4b, 0x13, 0x51, 0xd7, 0xfa, 0xaf, 0x98, 0xd6, 0xf1, 0x19, 0xc4,
	0x8d, 0x2c, 0xeb, 0x4c, 0x5e, 0x2f, 0x9b, 0x4c, 0x6d, 0x92, 0x0e, 0x75, 0xd3, 0xdb, 0x6b, 0x9f,
	0xb4, 0x84, 0x4f, 0x00, 0xbe, 0xcb, 0x52, 0x71, 0x5b, 0xe0, 0x51, 0x41, 0x44, 0x0a, 0xd9, 0x0f,
	0x21, 0x10, 0x45, 0xd1, 0x72, 0x95, 0xf8, 0xda, 0xea, 0xb0, 0x3d, 0x33, 0x7a, 0xc5, 0xb7, 0x6b,
	0xfd, 0x93, 0xc0, 0xea, 0x96, 0xe1, 0x23, 0x08, 0x8b, 0xb2, 0xe2, 0x4b, 0xdd, 0x64, 0xa8, 0x8d,
	0x98, 0x05, 0x86, 0x4e, 0x73, 0x44, 0xf0, 0xf2, 0x4c, 0x65, 0x49, 0x97, 0x54, 0xc2, 0x46, 0xdb,
	0xed, 0x74, 0x65, 0x44, 0xaf, 0x12, 0x36, 0xda, 0x26, 0x6b, 0x37, 0x09, 0x58, 0xcd, 0x60, 0xec,
	0x43, 0xc7, 0x64, 0x19, 0x93, 0x64, 0xe0, 0xf0, 0x87, 0x0b, 0xb1, 0x4d, 0x65, 0xa1, 0x32, 0xb5,
	0x6b, 0xff, 0x91, 0x4b, 0xb4, 0x12, 0x75, 0x53, 0x71, 0xc5, 0x73, 0x8a, 0xa7, 0xcb, 0x0e, 0x02,
	0xde, 0x07, 0x9f, 0x4b, 0x29, 0x24, 0x05, 0xe2, 0x33, 0x4b, 0x8e, 0x66, 0xf5, 0x6e, 0x99, 0xd5,
	0xff, 0x63, 0xd6, 0xe7, 0x10, 0x6c, 0x68, 0x69, 0x94, 0x41, 0x6f, 0xd4, 0xa3, 0xfc, 0xed, 0x1e,
	0xd9, 0xde, 0xba, 0x3d, 0x10, 0xdd, 0x43, 0x51, 0x65, 0xeb, 0x96, 0x12, 0xd1, 0x3d, 0x10, 0xf9,
	0xdf, 0x48, 0x7a, 0x87, 0x48, 0x42, 0xf0, 0xd3, 0xba, 0x51, 0xd7, 0x2f, 0x53, 0x08, 0xf7, 0x37,
	0x80, 0x5d, 0xf0, 0x66, 0xf3, 0x59, 0xda, 0xbf, 0x83, 0x3d, 0x08, 0xc7, 0x6c, 0x72, 0x31, 0xfd,
	0x9a, 0xf6, 0x1d, 0x43, 0x58, 0xba, 0xf8, 0x3c, 0x67, 0x69, 0xdf, 0x45, 0x80, 0x80, 0xa5, 0x97,
	0x73, 0x6d, 0x74, 0x0c, 0x9e, 0x8c, 0x67, 0x93, 0xf4, 0x63, 0xdf, 0x1b, 0xfd, 0x74, 0x20, 0x7a,
	0xaf, 0xb7, 0x74, 0x29, 0xae, 0xb8, 0xc4, 0x17, 0xd0, 0x65, 0x7c, 0x5d, 0xb6, 0x4a, 0xe3, 0x13,
	0x9a, 0xf3, 0xe6, 0xd8, 0x07, 0xc7, 0x63, 0xe3, 0x2b, 0x80, 0x0f, 0x5c, 0xd9, 0xd5, 0xb4, 0x78,
	0x6c, 0x0d, 0xee, 0x12, 0x39, 0xdc, 0xf2, 0x1b, 0x07, 0x5f, 0x43, 0x6c, 0xf7, 0xb7, 0x50, 0x92,
	0x67, 0x35, 0xde, 0x3b, 0x2a, 0xb1, 0xc6, 0x00, 0xec, 0x63, 0x66, 0xb2, 0x33, 0xe7, 0x5b, 0x40,
	0x1f, 0xdc, 0xdb, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xca, 0x2e, 0x7c, 0xd3, 0x7d, 0x03, 0x00,
	0x00,
}
