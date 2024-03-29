// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: file-copy/file-copy.proto

package file_copy

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

// FileCopyClient is the client API for FileCopy service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileCopyClient interface {
	Write(ctx context.Context, in *WriteArgs, opts ...grpc.CallOption) (*WriteResponse, error)
}

type fileCopyClient struct {
	cc grpc.ClientConnInterface
}

func NewFileCopyClient(cc grpc.ClientConnInterface) FileCopyClient {
	return &fileCopyClient{cc}
}

func (c *fileCopyClient) Write(ctx context.Context, in *WriteArgs, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, "/filecopy.FileCopy/Write", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileCopyServer is the server API for FileCopy service.
// All implementations must embed UnimplementedFileCopyServer
// for forward compatibility
type FileCopyServer interface {
	Write(context.Context, *WriteArgs) (*WriteResponse, error)
	mustEmbedUnimplementedFileCopyServer()
}

// UnimplementedFileCopyServer must be embedded to have forward compatible implementations.
type UnimplementedFileCopyServer struct {
}

func (UnimplementedFileCopyServer) Write(context.Context, *WriteArgs) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedFileCopyServer) mustEmbedUnimplementedFileCopyServer() {}

// UnsafeFileCopyServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileCopyServer will
// result in compilation errors.
type UnsafeFileCopyServer interface {
	mustEmbedUnimplementedFileCopyServer()
}

func RegisterFileCopyServer(s grpc.ServiceRegistrar, srv FileCopyServer) {
	s.RegisterService(&FileCopy_ServiceDesc, srv)
}

func _FileCopy_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteArgs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileCopyServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filecopy.FileCopy/Write",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileCopyServer).Write(ctx, req.(*WriteArgs))
	}
	return interceptor(ctx, in, info, handler)
}

// FileCopy_ServiceDesc is the grpc.ServiceDesc for FileCopy service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileCopy_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "filecopy.FileCopy",
	HandlerType: (*FileCopyServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Write",
			Handler:    _FileCopy_Write_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "file-copy/file-copy.proto",
}
