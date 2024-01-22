// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: bidstream.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	BidirectionalService_BidirectionalStream_FullMethodName = "/bidstream.BidirectionalService/BidirectionalStream"
)

// BidirectionalServiceClient is the client API for BidirectionalService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BidirectionalServiceClient interface {
	BidirectionalStream(ctx context.Context, opts ...grpc.CallOption) (BidirectionalService_BidirectionalStreamClient, error)
}

type bidirectionalServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBidirectionalServiceClient(cc grpc.ClientConnInterface) BidirectionalServiceClient {
	return &bidirectionalServiceClient{cc}
}

func (c *bidirectionalServiceClient) BidirectionalStream(ctx context.Context, opts ...grpc.CallOption) (BidirectionalService_BidirectionalStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &BidirectionalService_ServiceDesc.Streams[0], BidirectionalService_BidirectionalStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &bidirectionalServiceBidirectionalStreamClient{stream}
	return x, nil
}

type BidirectionalService_BidirectionalStreamClient interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ClientStream
}

type bidirectionalServiceBidirectionalStreamClient struct {
	grpc.ClientStream
}

func (x *bidirectionalServiceBidirectionalStreamClient) Send(m *Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *bidirectionalServiceBidirectionalStreamClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BidirectionalServiceServer is the server API for BidirectionalService service.
// All implementations must embed UnimplementedBidirectionalServiceServer
// for forward compatibility
type BidirectionalServiceServer interface {
	BidirectionalStream(BidirectionalService_BidirectionalStreamServer) error
	mustEmbedUnimplementedBidirectionalServiceServer()
}

// UnimplementedBidirectionalServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBidirectionalServiceServer struct {
}

func (UnimplementedBidirectionalServiceServer) BidirectionalStream(BidirectionalService_BidirectionalStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method BidirectionalStream not implemented")
}
func (UnimplementedBidirectionalServiceServer) mustEmbedUnimplementedBidirectionalServiceServer() {}

// UnsafeBidirectionalServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BidirectionalServiceServer will
// result in compilation errors.
type UnsafeBidirectionalServiceServer interface {
	mustEmbedUnimplementedBidirectionalServiceServer()
}

func RegisterBidirectionalServiceServer(s grpc.ServiceRegistrar, srv BidirectionalServiceServer) {
	s.RegisterService(&BidirectionalService_ServiceDesc, srv)
}

func _BidirectionalService_BidirectionalStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BidirectionalServiceServer).BidirectionalStream(&bidirectionalServiceBidirectionalStreamServer{stream})
}

type BidirectionalService_BidirectionalStreamServer interface {
	Send(*Message) error
	Recv() (*Message, error)
	grpc.ServerStream
}

type bidirectionalServiceBidirectionalStreamServer struct {
	grpc.ServerStream
}

func (x *bidirectionalServiceBidirectionalStreamServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func (x *bidirectionalServiceBidirectionalStreamServer) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BidirectionalService_ServiceDesc is the grpc.ServiceDesc for BidirectionalService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BidirectionalService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bidstream.BidirectionalService",
	HandlerType: (*BidirectionalServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BidirectionalStream",
			Handler:       _BidirectionalService_BidirectionalStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "bidstream.proto",
}
