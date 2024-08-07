// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: kuiper.proto

package api

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

// KuiperClient is the client API for Kuiper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KuiperClient interface {
	PutStandaloneConfig(ctx context.Context, in *NewStandaloneConfig, opts ...grpc.CallOption) (*StandaloneConfig, error)
	GetStandaloneConfig(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*StandaloneConfig, error)
	ListStandaloneConfig(ctx context.Context, in *ListStandaloneConfigReq, opts ...grpc.CallOption) (*ListStandaloneConfigResp, error)
	DeleteStandaloneConfig(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*StandaloneConfig, error)
	PlaceStandaloneConfig(ctx context.Context, in *PlaceReq, opts ...grpc.CallOption) (*PlaceResp, error)
	ListPlacementTaskByStandaloneConfig(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*ListPlacementTaskResp, error)
	DiffStandaloneConfig(ctx context.Context, in *DiffReq, opts ...grpc.CallOption) (*DiffStandaloneConfigResp, error)
	PutConfigGroup(ctx context.Context, in *NewConfigGroup, opts ...grpc.CallOption) (*ConfigGroup, error)
	GetConfigGroup(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*ConfigGroup, error)
	ListConfigGroup(ctx context.Context, in *ListConfigGroupReq, opts ...grpc.CallOption) (*ListConfigGroupResp, error)
	DeleteConfigGroup(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*ConfigGroup, error)
	PlaceConfigGroup(ctx context.Context, in *PlaceReq, opts ...grpc.CallOption) (*PlaceResp, error)
	ListPlacementTaskByConfigGroup(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*ListPlacementTaskResp, error)
	DiffConfigGroup(ctx context.Context, in *DiffReq, opts ...grpc.CallOption) (*DiffConfigGroupResp, error)
}

type kuiperClient struct {
	cc grpc.ClientConnInterface
}

func NewKuiperClient(cc grpc.ClientConnInterface) KuiperClient {
	return &kuiperClient{cc}
}

func (c *kuiperClient) PutStandaloneConfig(ctx context.Context, in *NewStandaloneConfig, opts ...grpc.CallOption) (*StandaloneConfig, error) {
	out := new(StandaloneConfig)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/PutStandaloneConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) GetStandaloneConfig(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*StandaloneConfig, error) {
	out := new(StandaloneConfig)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/GetStandaloneConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) ListStandaloneConfig(ctx context.Context, in *ListStandaloneConfigReq, opts ...grpc.CallOption) (*ListStandaloneConfigResp, error) {
	out := new(ListStandaloneConfigResp)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/ListStandaloneConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) DeleteStandaloneConfig(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*StandaloneConfig, error) {
	out := new(StandaloneConfig)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/DeleteStandaloneConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) PlaceStandaloneConfig(ctx context.Context, in *PlaceReq, opts ...grpc.CallOption) (*PlaceResp, error) {
	out := new(PlaceResp)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/PlaceStandaloneConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) ListPlacementTaskByStandaloneConfig(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*ListPlacementTaskResp, error) {
	out := new(ListPlacementTaskResp)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/ListPlacementTaskByStandaloneConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) DiffStandaloneConfig(ctx context.Context, in *DiffReq, opts ...grpc.CallOption) (*DiffStandaloneConfigResp, error) {
	out := new(DiffStandaloneConfigResp)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/DiffStandaloneConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) PutConfigGroup(ctx context.Context, in *NewConfigGroup, opts ...grpc.CallOption) (*ConfigGroup, error) {
	out := new(ConfigGroup)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/PutConfigGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) GetConfigGroup(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*ConfigGroup, error) {
	out := new(ConfigGroup)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/GetConfigGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) ListConfigGroup(ctx context.Context, in *ListConfigGroupReq, opts ...grpc.CallOption) (*ListConfigGroupResp, error) {
	out := new(ListConfigGroupResp)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/ListConfigGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) DeleteConfigGroup(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*ConfigGroup, error) {
	out := new(ConfigGroup)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/DeleteConfigGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) PlaceConfigGroup(ctx context.Context, in *PlaceReq, opts ...grpc.CallOption) (*PlaceResp, error) {
	out := new(PlaceResp)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/PlaceConfigGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) ListPlacementTaskByConfigGroup(ctx context.Context, in *ConfigId, opts ...grpc.CallOption) (*ListPlacementTaskResp, error) {
	out := new(ListPlacementTaskResp)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/ListPlacementTaskByConfigGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kuiperClient) DiffConfigGroup(ctx context.Context, in *DiffReq, opts ...grpc.CallOption) (*DiffConfigGroupResp, error) {
	out := new(DiffConfigGroupResp)
	err := c.cc.Invoke(ctx, "/proto.Kuiper/DiffConfigGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KuiperServer is the server API for Kuiper service.
// All implementations must embed UnimplementedKuiperServer
// for forward compatibility
type KuiperServer interface {
	PutStandaloneConfig(context.Context, *NewStandaloneConfig) (*StandaloneConfig, error)
	GetStandaloneConfig(context.Context, *ConfigId) (*StandaloneConfig, error)
	ListStandaloneConfig(context.Context, *ListStandaloneConfigReq) (*ListStandaloneConfigResp, error)
	DeleteStandaloneConfig(context.Context, *ConfigId) (*StandaloneConfig, error)
	PlaceStandaloneConfig(context.Context, *PlaceReq) (*PlaceResp, error)
	ListPlacementTaskByStandaloneConfig(context.Context, *ConfigId) (*ListPlacementTaskResp, error)
	DiffStandaloneConfig(context.Context, *DiffReq) (*DiffStandaloneConfigResp, error)
	PutConfigGroup(context.Context, *NewConfigGroup) (*ConfigGroup, error)
	GetConfigGroup(context.Context, *ConfigId) (*ConfigGroup, error)
	ListConfigGroup(context.Context, *ListConfigGroupReq) (*ListConfigGroupResp, error)
	DeleteConfigGroup(context.Context, *ConfigId) (*ConfigGroup, error)
	PlaceConfigGroup(context.Context, *PlaceReq) (*PlaceResp, error)
	ListPlacementTaskByConfigGroup(context.Context, *ConfigId) (*ListPlacementTaskResp, error)
	DiffConfigGroup(context.Context, *DiffReq) (*DiffConfigGroupResp, error)
	mustEmbedUnimplementedKuiperServer()
}

// UnimplementedKuiperServer must be embedded to have forward compatible implementations.
type UnimplementedKuiperServer struct {
}

func (UnimplementedKuiperServer) PutStandaloneConfig(context.Context, *NewStandaloneConfig) (*StandaloneConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutStandaloneConfig not implemented")
}
func (UnimplementedKuiperServer) GetStandaloneConfig(context.Context, *ConfigId) (*StandaloneConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStandaloneConfig not implemented")
}
func (UnimplementedKuiperServer) ListStandaloneConfig(context.Context, *ListStandaloneConfigReq) (*ListStandaloneConfigResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListStandaloneConfig not implemented")
}
func (UnimplementedKuiperServer) DeleteStandaloneConfig(context.Context, *ConfigId) (*StandaloneConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStandaloneConfig not implemented")
}
func (UnimplementedKuiperServer) PlaceStandaloneConfig(context.Context, *PlaceReq) (*PlaceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceStandaloneConfig not implemented")
}
func (UnimplementedKuiperServer) ListPlacementTaskByStandaloneConfig(context.Context, *ConfigId) (*ListPlacementTaskResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPlacementTaskByStandaloneConfig not implemented")
}
func (UnimplementedKuiperServer) DiffStandaloneConfig(context.Context, *DiffReq) (*DiffStandaloneConfigResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiffStandaloneConfig not implemented")
}
func (UnimplementedKuiperServer) PutConfigGroup(context.Context, *NewConfigGroup) (*ConfigGroup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PutConfigGroup not implemented")
}
func (UnimplementedKuiperServer) GetConfigGroup(context.Context, *ConfigId) (*ConfigGroup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfigGroup not implemented")
}
func (UnimplementedKuiperServer) ListConfigGroup(context.Context, *ListConfigGroupReq) (*ListConfigGroupResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListConfigGroup not implemented")
}
func (UnimplementedKuiperServer) DeleteConfigGroup(context.Context, *ConfigId) (*ConfigGroup, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteConfigGroup not implemented")
}
func (UnimplementedKuiperServer) PlaceConfigGroup(context.Context, *PlaceReq) (*PlaceResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlaceConfigGroup not implemented")
}
func (UnimplementedKuiperServer) ListPlacementTaskByConfigGroup(context.Context, *ConfigId) (*ListPlacementTaskResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPlacementTaskByConfigGroup not implemented")
}
func (UnimplementedKuiperServer) DiffConfigGroup(context.Context, *DiffReq) (*DiffConfigGroupResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DiffConfigGroup not implemented")
}
func (UnimplementedKuiperServer) mustEmbedUnimplementedKuiperServer() {}

// UnsafeKuiperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KuiperServer will
// result in compilation errors.
type UnsafeKuiperServer interface {
	mustEmbedUnimplementedKuiperServer()
}

func RegisterKuiperServer(s grpc.ServiceRegistrar, srv KuiperServer) {
	s.RegisterService(&Kuiper_ServiceDesc, srv)
}

func _Kuiper_PutStandaloneConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewStandaloneConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).PutStandaloneConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/PutStandaloneConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).PutStandaloneConfig(ctx, req.(*NewStandaloneConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_GetStandaloneConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).GetStandaloneConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/GetStandaloneConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).GetStandaloneConfig(ctx, req.(*ConfigId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_ListStandaloneConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStandaloneConfigReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).ListStandaloneConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/ListStandaloneConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).ListStandaloneConfig(ctx, req.(*ListStandaloneConfigReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_DeleteStandaloneConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).DeleteStandaloneConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/DeleteStandaloneConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).DeleteStandaloneConfig(ctx, req.(*ConfigId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_PlaceStandaloneConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).PlaceStandaloneConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/PlaceStandaloneConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).PlaceStandaloneConfig(ctx, req.(*PlaceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_ListPlacementTaskByStandaloneConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).ListPlacementTaskByStandaloneConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/ListPlacementTaskByStandaloneConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).ListPlacementTaskByStandaloneConfig(ctx, req.(*ConfigId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_DiffStandaloneConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiffReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).DiffStandaloneConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/DiffStandaloneConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).DiffStandaloneConfig(ctx, req.(*DiffReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_PutConfigGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewConfigGroup)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).PutConfigGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/PutConfigGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).PutConfigGroup(ctx, req.(*NewConfigGroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_GetConfigGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).GetConfigGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/GetConfigGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).GetConfigGroup(ctx, req.(*ConfigId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_ListConfigGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListConfigGroupReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).ListConfigGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/ListConfigGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).ListConfigGroup(ctx, req.(*ListConfigGroupReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_DeleteConfigGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).DeleteConfigGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/DeleteConfigGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).DeleteConfigGroup(ctx, req.(*ConfigId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_PlaceConfigGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlaceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).PlaceConfigGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/PlaceConfigGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).PlaceConfigGroup(ctx, req.(*PlaceReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_ListPlacementTaskByConfigGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).ListPlacementTaskByConfigGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/ListPlacementTaskByConfigGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).ListPlacementTaskByConfigGroup(ctx, req.(*ConfigId))
	}
	return interceptor(ctx, in, info, handler)
}

func _Kuiper_DiffConfigGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DiffReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuiperServer).DiffConfigGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Kuiper/DiffConfigGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuiperServer).DiffConfigGroup(ctx, req.(*DiffReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Kuiper_ServiceDesc is the grpc.ServiceDesc for Kuiper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Kuiper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Kuiper",
	HandlerType: (*KuiperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PutStandaloneConfig",
			Handler:    _Kuiper_PutStandaloneConfig_Handler,
		},
		{
			MethodName: "GetStandaloneConfig",
			Handler:    _Kuiper_GetStandaloneConfig_Handler,
		},
		{
			MethodName: "ListStandaloneConfig",
			Handler:    _Kuiper_ListStandaloneConfig_Handler,
		},
		{
			MethodName: "DeleteStandaloneConfig",
			Handler:    _Kuiper_DeleteStandaloneConfig_Handler,
		},
		{
			MethodName: "PlaceStandaloneConfig",
			Handler:    _Kuiper_PlaceStandaloneConfig_Handler,
		},
		{
			MethodName: "ListPlacementTaskByStandaloneConfig",
			Handler:    _Kuiper_ListPlacementTaskByStandaloneConfig_Handler,
		},
		{
			MethodName: "DiffStandaloneConfig",
			Handler:    _Kuiper_DiffStandaloneConfig_Handler,
		},
		{
			MethodName: "PutConfigGroup",
			Handler:    _Kuiper_PutConfigGroup_Handler,
		},
		{
			MethodName: "GetConfigGroup",
			Handler:    _Kuiper_GetConfigGroup_Handler,
		},
		{
			MethodName: "ListConfigGroup",
			Handler:    _Kuiper_ListConfigGroup_Handler,
		},
		{
			MethodName: "DeleteConfigGroup",
			Handler:    _Kuiper_DeleteConfigGroup_Handler,
		},
		{
			MethodName: "PlaceConfigGroup",
			Handler:    _Kuiper_PlaceConfigGroup_Handler,
		},
		{
			MethodName: "ListPlacementTaskByConfigGroup",
			Handler:    _Kuiper_ListPlacementTaskByConfigGroup_Handler,
		},
		{
			MethodName: "DiffConfigGroup",
			Handler:    _Kuiper_DiffConfigGroup_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kuiper.proto",
}
