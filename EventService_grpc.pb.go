// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: odpf/raccoon/EventService.proto

package raccoon

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

// EventServiceClient is the client API for EventService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventServiceClient interface {
	SendEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*EventResponse, error)
}

type eventServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEventServiceClient(cc grpc.ClientConnInterface) EventServiceClient {
	return &eventServiceClient{cc}
}

func (c *eventServiceClient) SendEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/odpf.raccoon.EventService/SendEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventServiceServer is the server API for EventService service.
// All implementations must embed UnimplementedEventServiceServer
// for forward compatibility
type EventServiceServer interface {
	SendEvent(context.Context, *EventRequest) (*EventResponse, error)
	mustEmbedUnimplementedEventServiceServer()
}

// UnimplementedEventServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEventServiceServer struct {
}

func (UnimplementedEventServiceServer) SendEvent(context.Context, *EventRequest) (*EventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEvent not implemented")
}
func (UnimplementedEventServiceServer) mustEmbedUnimplementedEventServiceServer() {}

// UnsafeEventServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventServiceServer will
// result in compilation errors.
type UnsafeEventServiceServer interface {
	mustEmbedUnimplementedEventServiceServer()
}

func RegisterEventServiceServer(s grpc.ServiceRegistrar, srv EventServiceServer) {
	s.RegisterService(&EventService_ServiceDesc, srv)
}

func _EventService_SendEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventServiceServer).SendEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/odpf.raccoon.EventService/SendEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventServiceServer).SendEvent(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EventService_ServiceDesc is the grpc.ServiceDesc for EventService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "odpf.raccoon.EventService",
	HandlerType: (*EventServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEvent",
			Handler:    _EventService_SendEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "odpf/raccoon/EventService.proto",
}
