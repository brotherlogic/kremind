// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: kremind.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	KremindService_AddReminder_FullMethodName   = "/kremind.KremindService/AddReminder"
	KremindService_ListReminders_FullMethodName = "/kremind.KremindService/ListReminders"
)

// KremindServiceClient is the client API for KremindService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KremindServiceClient interface {
	AddReminder(ctx context.Context, in *AddReminderRequest, opts ...grpc.CallOption) (*AddReminderResponse, error)
	ListReminders(ctx context.Context, in *ListRemindersRequest, opts ...grpc.CallOption) (*ListRemindersResponse, error)
}

type kremindServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKremindServiceClient(cc grpc.ClientConnInterface) KremindServiceClient {
	return &kremindServiceClient{cc}
}

func (c *kremindServiceClient) AddReminder(ctx context.Context, in *AddReminderRequest, opts ...grpc.CallOption) (*AddReminderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddReminderResponse)
	err := c.cc.Invoke(ctx, KremindService_AddReminder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kremindServiceClient) ListReminders(ctx context.Context, in *ListRemindersRequest, opts ...grpc.CallOption) (*ListRemindersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListRemindersResponse)
	err := c.cc.Invoke(ctx, KremindService_ListReminders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KremindServiceServer is the server API for KremindService service.
// All implementations should embed UnimplementedKremindServiceServer
// for forward compatibility.
type KremindServiceServer interface {
	AddReminder(context.Context, *AddReminderRequest) (*AddReminderResponse, error)
	ListReminders(context.Context, *ListRemindersRequest) (*ListRemindersResponse, error)
}

// UnimplementedKremindServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedKremindServiceServer struct{}

func (UnimplementedKremindServiceServer) AddReminder(context.Context, *AddReminderRequest) (*AddReminderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddReminder not implemented")
}
func (UnimplementedKremindServiceServer) ListReminders(context.Context, *ListRemindersRequest) (*ListRemindersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListReminders not implemented")
}
func (UnimplementedKremindServiceServer) testEmbeddedByValue() {}

// UnsafeKremindServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KremindServiceServer will
// result in compilation errors.
type UnsafeKremindServiceServer interface {
	mustEmbedUnimplementedKremindServiceServer()
}

func RegisterKremindServiceServer(s grpc.ServiceRegistrar, srv KremindServiceServer) {
	// If the following call pancis, it indicates UnimplementedKremindServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&KremindService_ServiceDesc, srv)
}

func _KremindService_AddReminder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddReminderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KremindServiceServer).AddReminder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KremindService_AddReminder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KremindServiceServer).AddReminder(ctx, req.(*AddReminderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KremindService_ListReminders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRemindersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KremindServiceServer).ListReminders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KremindService_ListReminders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KremindServiceServer).ListReminders(ctx, req.(*ListRemindersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KremindService_ServiceDesc is the grpc.ServiceDesc for KremindService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KremindService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kremind.KremindService",
	HandlerType: (*KremindServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddReminder",
			Handler:    _KremindService_AddReminder_Handler,
		},
		{
			MethodName: "ListReminders",
			Handler:    _KremindService_ListReminders_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "kremind.proto",
}
