// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.17.3
// source: tcp/tcp_rpc.proto

package tcp

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

// TCPRPCClient is the client API for TCPRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TCPRPCClient interface {
	HelloCaller(ctx context.Context, in *HelloMsg, opts ...grpc.CallOption) (*OK, error)
	LoginCaller(ctx context.Context, in *UserMsg, opts ...grpc.CallOption) (*UserMsg, error)
	NicknameCaller(ctx context.Context, in *UserMsg, opts ...grpc.CallOption) (*OK, error)
	PhotoCaller(ctx context.Context, in *UserMsg, opts ...grpc.CallOption) (*OK, error)
}

type tCPRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewTCPRPCClient(cc grpc.ClientConnInterface) TCPRPCClient {
	return &tCPRPCClient{cc}
}

func (c *tCPRPCClient) HelloCaller(ctx context.Context, in *HelloMsg, opts ...grpc.CallOption) (*OK, error) {
	out := new(OK)
	err := c.cc.Invoke(ctx, "/tcp.TCPRPC/HelloCaller", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tCPRPCClient) LoginCaller(ctx context.Context, in *UserMsg, opts ...grpc.CallOption) (*UserMsg, error) {
	out := new(UserMsg)
	err := c.cc.Invoke(ctx, "/tcp.TCPRPC/LoginCaller", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tCPRPCClient) NicknameCaller(ctx context.Context, in *UserMsg, opts ...grpc.CallOption) (*OK, error) {
	out := new(OK)
	err := c.cc.Invoke(ctx, "/tcp.TCPRPC/NicknameCaller", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tCPRPCClient) PhotoCaller(ctx context.Context, in *UserMsg, opts ...grpc.CallOption) (*OK, error) {
	out := new(OK)
	err := c.cc.Invoke(ctx, "/tcp.TCPRPC/PhotoCaller", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TCPRPCServer is the server API for TCPRPC service.
// All implementations must embed UnimplementedTCPRPCServer
// for forward compatibility
type TCPRPCServer interface {
	HelloCaller(context.Context, *HelloMsg) (*OK, error)
	LoginCaller(context.Context, *UserMsg) (*UserMsg, error)
	NicknameCaller(context.Context, *UserMsg) (*OK, error)
	PhotoCaller(context.Context, *UserMsg) (*OK, error)
	mustEmbedUnimplementedTCPRPCServer()
}

// UnimplementedTCPRPCServer must be embedded to have forward compatible implementations.
type UnimplementedTCPRPCServer struct {
}

func (UnimplementedTCPRPCServer) HelloCaller(context.Context, *HelloMsg) (*OK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HelloCaller not implemented")
}
func (UnimplementedTCPRPCServer) LoginCaller(context.Context, *UserMsg) (*UserMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginCaller not implemented")
}
func (UnimplementedTCPRPCServer) NicknameCaller(context.Context, *UserMsg) (*OK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NicknameCaller not implemented")
}
func (UnimplementedTCPRPCServer) PhotoCaller(context.Context, *UserMsg) (*OK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PhotoCaller not implemented")
}
func (UnimplementedTCPRPCServer) mustEmbedUnimplementedTCPRPCServer() {}

// UnsafeTCPRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TCPRPCServer will
// result in compilation errors.
type UnsafeTCPRPCServer interface {
	mustEmbedUnimplementedTCPRPCServer()
}

func RegisterTCPRPCServer(s grpc.ServiceRegistrar, srv TCPRPCServer) {
	s.RegisterService(&TCPRPC_ServiceDesc, srv)
}

func _TCPRPC_HelloCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TCPRPCServer).HelloCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tcp.TCPRPC/HelloCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TCPRPCServer).HelloCaller(ctx, req.(*HelloMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _TCPRPC_LoginCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TCPRPCServer).LoginCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tcp.TCPRPC/LoginCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TCPRPCServer).LoginCaller(ctx, req.(*UserMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _TCPRPC_NicknameCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TCPRPCServer).NicknameCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tcp.TCPRPC/NicknameCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TCPRPCServer).NicknameCaller(ctx, req.(*UserMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _TCPRPC_PhotoCaller_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TCPRPCServer).PhotoCaller(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tcp.TCPRPC/PhotoCaller",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TCPRPCServer).PhotoCaller(ctx, req.(*UserMsg))
	}
	return interceptor(ctx, in, info, handler)
}

// TCPRPC_ServiceDesc is the grpc.ServiceDesc for TCPRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TCPRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tcp.TCPRPC",
	HandlerType: (*TCPRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HelloCaller",
			Handler:    _TCPRPC_HelloCaller_Handler,
		},
		{
			MethodName: "LoginCaller",
			Handler:    _TCPRPC_LoginCaller_Handler,
		},
		{
			MethodName: "NicknameCaller",
			Handler:    _TCPRPC_NicknameCaller_Handler,
		},
		{
			MethodName: "PhotoCaller",
			Handler:    _TCPRPC_PhotoCaller_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tcp/tcp_rpc.proto",
}