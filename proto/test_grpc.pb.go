// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: proto/test.proto

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

// GeometryServiceClient is the client API for GeometryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GeometryServiceClient interface {
	// методы, которые можно будет реализовать и использовать
	UpdateAgent(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	SendDataToAgent(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Empty, error)
}

type geometryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGeometryServiceClient(cc grpc.ClientConnInterface) GeometryServiceClient {
	return &geometryServiceClient{cc}
}

func (c *geometryServiceClient) UpdateAgent(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/test.GeometryService/Update_agent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *geometryServiceClient) SendDataToAgent(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/test.GeometryService/Send_data_to_agent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GeometryServiceServer is the server API for GeometryService service.
// All implementations must embed UnimplementedGeometryServiceServer
// for forward compatibility
type GeometryServiceServer interface {
	// методы, которые можно будет реализовать и использовать
	UpdateAgent(context.Context, *Empty) (*Empty, error)
	SendDataToAgent(context.Context, *Data) (*Empty, error)
	mustEmbedUnimplementedGeometryServiceServer()
}

// UnimplementedGeometryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGeometryServiceServer struct {
}

func (UnimplementedGeometryServiceServer) UpdateAgent(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAgent not implemented")
}
func (UnimplementedGeometryServiceServer) SendDataToAgent(context.Context, *Data) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendDataToAgent not implemented")
}
func (UnimplementedGeometryServiceServer) mustEmbedUnimplementedGeometryServiceServer() {}

// UnsafeGeometryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GeometryServiceServer will
// result in compilation errors.
type UnsafeGeometryServiceServer interface {
	mustEmbedUnimplementedGeometryServiceServer()
}

func RegisterGeometryServiceServer(s grpc.ServiceRegistrar, srv GeometryServiceServer) {
	s.RegisterService(&GeometryService_ServiceDesc, srv)
}

func _GeometryService_UpdateAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeometryServiceServer).UpdateAgent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.GeometryService/Update_agent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeometryServiceServer).UpdateAgent(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _GeometryService_SendDataToAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Data)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GeometryServiceServer).SendDataToAgent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/test.GeometryService/Send_data_to_agent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GeometryServiceServer).SendDataToAgent(ctx, req.(*Data))
	}
	return interceptor(ctx, in, info, handler)
}

// GeometryService_ServiceDesc is the grpc.ServiceDesc for GeometryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GeometryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "test.GeometryService",
	HandlerType: (*GeometryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update_agent",
			Handler:    _GeometryService_UpdateAgent_Handler,
		},
		{
			MethodName: "Send_data_to_agent",
			Handler:    _GeometryService_SendDataToAgent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/test.proto",
}