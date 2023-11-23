// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.15.5
// source: idl/MapleJuice.proto

package idl

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
	MapleJuiceService_ExecuteMaple_FullMethodName = "/maplejuice.MapleJuiceService/ExecuteMaple"
	MapleJuiceService_ExecuteJuice_FullMethodName = "/maplejuice.MapleJuiceService/ExecuteJuice"
)

// MapleJuiceServiceClient is the client API for MapleJuiceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MapleJuiceServiceClient interface {
	ExecuteMaple(ctx context.Context, in *MapleRequest, opts ...grpc.CallOption) (*MapleResponse, error)
	ExecuteJuice(ctx context.Context, in *JuiceRequest, opts ...grpc.CallOption) (*JuiceResponse, error)
}

type mapleJuiceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMapleJuiceServiceClient(cc grpc.ClientConnInterface) MapleJuiceServiceClient {
	return &mapleJuiceServiceClient{cc}
}

func (c *mapleJuiceServiceClient) ExecuteMaple(ctx context.Context, in *MapleRequest, opts ...grpc.CallOption) (*MapleResponse, error) {
	out := new(MapleResponse)
	err := c.cc.Invoke(ctx, MapleJuiceService_ExecuteMaple_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mapleJuiceServiceClient) ExecuteJuice(ctx context.Context, in *JuiceRequest, opts ...grpc.CallOption) (*JuiceResponse, error) {
	out := new(JuiceResponse)
	err := c.cc.Invoke(ctx, MapleJuiceService_ExecuteJuice_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MapleJuiceServiceServer is the server API for MapleJuiceService service.
// All implementations must embed UnimplementedMapleJuiceServiceServer
// for forward compatibility
type MapleJuiceServiceServer interface {
	ExecuteMaple(context.Context, *MapleRequest) (*MapleResponse, error)
	ExecuteJuice(context.Context, *JuiceRequest) (*JuiceResponse, error)
	mustEmbedUnimplementedMapleJuiceServiceServer()
}

// UnimplementedMapleJuiceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMapleJuiceServiceServer struct {
}

func (UnimplementedMapleJuiceServiceServer) ExecuteMaple(context.Context, *MapleRequest) (*MapleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteMaple not implemented")
}
func (UnimplementedMapleJuiceServiceServer) ExecuteJuice(context.Context, *JuiceRequest) (*JuiceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteJuice not implemented")
}
func (UnimplementedMapleJuiceServiceServer) mustEmbedUnimplementedMapleJuiceServiceServer() {}

// UnsafeMapleJuiceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MapleJuiceServiceServer will
// result in compilation errors.
type UnsafeMapleJuiceServiceServer interface {
	mustEmbedUnimplementedMapleJuiceServiceServer()
}

func RegisterMapleJuiceServiceServer(s grpc.ServiceRegistrar, srv MapleJuiceServiceServer) {
	s.RegisterService(&MapleJuiceService_ServiceDesc, srv)
}

func _MapleJuiceService_ExecuteMaple_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MapleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapleJuiceServiceServer).ExecuteMaple(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MapleJuiceService_ExecuteMaple_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapleJuiceServiceServer).ExecuteMaple(ctx, req.(*MapleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MapleJuiceService_ExecuteJuice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JuiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MapleJuiceServiceServer).ExecuteJuice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MapleJuiceService_ExecuteJuice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MapleJuiceServiceServer).ExecuteJuice(ctx, req.(*JuiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MapleJuiceService_ServiceDesc is the grpc.ServiceDesc for MapleJuiceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MapleJuiceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "maplejuice.MapleJuiceService",
	HandlerType: (*MapleJuiceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExecuteMaple",
			Handler:    _MapleJuiceService_ExecuteMaple_Handler,
		},
		{
			MethodName: "ExecuteJuice",
			Handler:    _MapleJuiceService_ExecuteJuice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "idl/MapleJuice.proto",
}
