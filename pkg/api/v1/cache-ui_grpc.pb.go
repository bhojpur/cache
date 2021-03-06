// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// CacheUIClient is the client API for CacheUI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CacheUIClient interface {
	// ListEngineSpecs returns a list of Cache Engine(s) that can be started through the UI.
	ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (CacheUI_ListEngineSpecsClient, error)
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error)
}

type cacheUIClient struct {
	cc grpc.ClientConnInterface
}

func NewCacheUIClient(cc grpc.ClientConnInterface) CacheUIClient {
	return &cacheUIClient{cc}
}

func (c *cacheUIClient) ListEngineSpecs(ctx context.Context, in *ListEngineSpecsRequest, opts ...grpc.CallOption) (CacheUI_ListEngineSpecsClient, error) {
	stream, err := c.cc.NewStream(ctx, &CacheUI_ServiceDesc.Streams[0], "/v1.CacheUI/ListEngineSpecs", opts...)
	if err != nil {
		return nil, err
	}
	x := &cacheUIListEngineSpecsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CacheUI_ListEngineSpecsClient interface {
	Recv() (*ListEngineSpecsResponse, error)
	grpc.ClientStream
}

type cacheUIListEngineSpecsClient struct {
	grpc.ClientStream
}

func (x *cacheUIListEngineSpecsClient) Recv() (*ListEngineSpecsResponse, error) {
	m := new(ListEngineSpecsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cacheUIClient) IsReadOnly(ctx context.Context, in *IsReadOnlyRequest, opts ...grpc.CallOption) (*IsReadOnlyResponse, error) {
	out := new(IsReadOnlyResponse)
	err := c.cc.Invoke(ctx, "/v1.CacheUI/IsReadOnly", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CacheUIServer is the server API for CacheUI service.
// All implementations must embed UnimplementedCacheUIServer
// for forward compatibility
type CacheUIServer interface {
	// ListEngineSpecs returns a list of Cache Engine(s) that can be started through the UI.
	ListEngineSpecs(*ListEngineSpecsRequest, CacheUI_ListEngineSpecsServer) error
	// IsReadOnly returns true if the UI is readonly.
	IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error)
	mustEmbedUnimplementedCacheUIServer()
}

// UnimplementedCacheUIServer must be embedded to have forward compatible implementations.
type UnimplementedCacheUIServer struct {
}

func (UnimplementedCacheUIServer) ListEngineSpecs(*ListEngineSpecsRequest, CacheUI_ListEngineSpecsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListEngineSpecs not implemented")
}
func (UnimplementedCacheUIServer) IsReadOnly(context.Context, *IsReadOnlyRequest) (*IsReadOnlyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsReadOnly not implemented")
}
func (UnimplementedCacheUIServer) mustEmbedUnimplementedCacheUIServer() {}

// UnsafeCacheUIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CacheUIServer will
// result in compilation errors.
type UnsafeCacheUIServer interface {
	mustEmbedUnimplementedCacheUIServer()
}

func RegisterCacheUIServer(s grpc.ServiceRegistrar, srv CacheUIServer) {
	s.RegisterService(&CacheUI_ServiceDesc, srv)
}

func _CacheUI_ListEngineSpecs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListEngineSpecsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CacheUIServer).ListEngineSpecs(m, &cacheUIListEngineSpecsServer{stream})
}

type CacheUI_ListEngineSpecsServer interface {
	Send(*ListEngineSpecsResponse) error
	grpc.ServerStream
}

type cacheUIListEngineSpecsServer struct {
	grpc.ServerStream
}

func (x *cacheUIListEngineSpecsServer) Send(m *ListEngineSpecsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _CacheUI_IsReadOnly_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IsReadOnlyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacheUIServer).IsReadOnly(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/v1.CacheUI/IsReadOnly",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacheUIServer).IsReadOnly(ctx, req.(*IsReadOnlyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CacheUI_ServiceDesc is the grpc.ServiceDesc for CacheUI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CacheUI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.CacheUI",
	HandlerType: (*CacheUIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IsReadOnly",
			Handler:    _CacheUI_IsReadOnly_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListEngineSpecs",
			Handler:       _CacheUI_ListEngineSpecs_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "cache-ui.proto",
}
