// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.17.3
// source: registry.proto

package __

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

// RegistryClient is the client API for Registry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RegistryClient interface {
	Register(ctx context.Context, in *RegInfo, opts ...grpc.CallOption) (*RegAnswer, error)
	StartGroup(ctx context.Context, in *RequestData, opts ...grpc.CallOption) (*Group, error)
	Ready(ctx context.Context, in *RequestData, opts ...grpc.CallOption) (*Group, error)
	CloseGroup(ctx context.Context, in *RequestData, opts ...grpc.CallOption) (*Group, error)
	GetStatus(ctx context.Context, in *MulticastId, opts ...grpc.CallOption) (*Group, error)
}

type registryClient struct {
	cc grpc.ClientConnInterface
}

func NewRegistryClient(cc grpc.ClientConnInterface) RegistryClient {
	return &registryClient{cc}
}

func (c *registryClient) Register(ctx context.Context, in *RegInfo, opts ...grpc.CallOption) (*RegAnswer, error) {
	out := new(RegAnswer)
	err := c.cc.Invoke(ctx, "/serviceregistry.Registry/register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) StartGroup(ctx context.Context, in *RequestData, opts ...grpc.CallOption) (*Group, error) {
	out := new(Group)
	err := c.cc.Invoke(ctx, "/serviceregistry.Registry/startGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) Ready(ctx context.Context, in *RequestData, opts ...grpc.CallOption) (*Group, error) {
	out := new(Group)
	err := c.cc.Invoke(ctx, "/serviceregistry.Registry/ready", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) CloseGroup(ctx context.Context, in *RequestData, opts ...grpc.CallOption) (*Group, error) {
	out := new(Group)
	err := c.cc.Invoke(ctx, "/serviceregistry.Registry/closeGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *registryClient) GetStatus(ctx context.Context, in *MulticastId, opts ...grpc.CallOption) (*Group, error) {
	out := new(Group)
	err := c.cc.Invoke(ctx, "/serviceregistry.Registry/getStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RegistryServer is the server API for Registry service.
// All implementations must embed UnimplementedRegistryServer
// for forward compatibility
type RegistryServer interface {
	Register(context.Context, *RegInfo) (*RegAnswer, error)
	StartGroup(context.Context, *RequestData) (*Group, error)
	Ready(context.Context, *RequestData) (*Group, error)
	CloseGroup(context.Context, *RequestData) (*Group, error)
	GetStatus(context.Context, *MulticastId) (*Group, error)
	mustEmbedUnimplementedRegistryServer()
}

// UnimplementedRegistryServer must be embedded to have forward compatible implementations.
type UnimplementedRegistryServer struct {
}

func (UnimplementedRegistryServer) Register(context.Context, *RegInfo) (*RegAnswer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRegistryServer) StartGroup(context.Context, *RequestData) (*Group, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartGroup not implemented")
}
func (UnimplementedRegistryServer) Ready(context.Context, *RequestData) (*Group, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ready not implemented")
}
func (UnimplementedRegistryServer) CloseGroup(context.Context, *RequestData) (*Group, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseGroup not implemented")
}
func (UnimplementedRegistryServer) GetStatus(context.Context, *MulticastId) (*Group, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedRegistryServer) mustEmbedUnimplementedRegistryServer() {}

// UnsafeRegistryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RegistryServer will
// result in compilation errors.
type UnsafeRegistryServer interface {
	mustEmbedUnimplementedRegistryServer()
}

func RegisterRegistryServer(s grpc.ServiceRegistrar, srv RegistryServer) {
	s.RegisterService(&Registry_ServiceDesc, srv)
}

func _Registry_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serviceregistry.Registry/register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).Register(ctx, req.(*RegInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_StartGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).StartGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serviceregistry.Registry/startGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).StartGroup(ctx, req.(*RequestData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_Ready_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).Ready(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serviceregistry.Registry/ready",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).Ready(ctx, req.(*RequestData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_CloseGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestData)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).CloseGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serviceregistry.Registry/closeGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).CloseGroup(ctx, req.(*RequestData))
	}
	return interceptor(ctx, in, info, handler)
}

func _Registry_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MulticastId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RegistryServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serviceregistry.Registry/getStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RegistryServer).GetStatus(ctx, req.(*MulticastId))
	}
	return interceptor(ctx, in, info, handler)
}

// Registry_ServiceDesc is the grpc.ServiceDesc for Registry service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Registry_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "serviceregistry.Registry",
	HandlerType: (*RegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "register",
			Handler:    _Registry_Register_Handler,
		},
		{
			MethodName: "startGroup",
			Handler:    _Registry_StartGroup_Handler,
		},
		{
			MethodName: "ready",
			Handler:    _Registry_Ready_Handler,
		},
		{
			MethodName: "closeGroup",
			Handler:    _Registry_CloseGroup_Handler,
		},
		{
			MethodName: "getStatus",
			Handler:    _Registry_GetStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "registry.proto",
}
