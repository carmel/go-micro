//
//
// Tencent is pleased to support the open source community by making tRPC available.
//
// Copyright (C) 2023 THL A29 Limited, a Tencent company.
// All rights reserved.
//
// If you have downloaded a copy of the tRPC source code from Tencent,
// please note that tRPC source code is licensed under the  Apache 2.0 License,
// A copy of the Apache 2.0 License is included in this file.
//
//

// Code generated by trpc-go/trpc-cmdline v2.0.17. DO NOT EDIT.
// source: helloworld.proto

package proto

import (
	"context"
	"errors"
	"fmt"

	_ "go-micro"
	"go-micro/client"
	"go-micro/codec"
	_ "go-micro/http"
	"go-micro/server"
	"go-micro/stream"
)

// START ======================================= Server Service Definition ======================================= START

// TestStreamService defines service
type TestStreamService interface {
	// ClientStream Defined Client-side streaming RPC
	//  Add stream in front of HelloReq
	ClientStream(TestStream_ClientStreamServer) error
	// ServerStream Defined Server-side streaming RPC
	//  Add stream in front of HelloRsp
	ServerStream(*HelloReq, TestStream_ServerStreamServer) error
	// BidirectionalStream Bidirectional streaming RPC
	//  Add stream in front of HelloReq and HelloRsp
	BidirectionalStream(TestStream_BidirectionalStreamServer) error
}

func TestStreamService_ClientStream_Handler(srv interface{}, stream server.Stream) error {
	return srv.(TestStreamService).ClientStream(&testStreamClientStreamServer{stream})
}

type TestStream_ClientStreamServer interface {
	SendAndClose(*HelloRsp) error
	Recv() (*HelloReq, error)
	server.Stream
}

type testStreamClientStreamServer struct {
	server.Stream
}

func (x *testStreamClientStreamServer) SendAndClose(m *HelloRsp) error {
	return x.Stream.SendMsg(m)
}

func (x *testStreamClientStreamServer) Recv() (*HelloReq, error) {
	m := new(HelloReq)
	if err := x.Stream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func TestStreamService_ServerStream_Handler(srv interface{}, stream server.Stream) error {
	m := new(HelloReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TestStreamService).ServerStream(m, &testStreamServerStreamServer{stream})
}

type TestStream_ServerStreamServer interface {
	Send(*HelloRsp) error
	server.Stream
}

type testStreamServerStreamServer struct {
	server.Stream
}

func (x *testStreamServerStreamServer) Send(m *HelloRsp) error {
	return x.Stream.SendMsg(m)
}

func TestStreamService_BidirectionalStream_Handler(srv interface{}, stream server.Stream) error {
	return srv.(TestStreamService).BidirectionalStream(&testStreamBidirectionalStreamServer{stream})
}

type TestStream_BidirectionalStreamServer interface {
	Send(*HelloRsp) error
	Recv() (*HelloReq, error)
	server.Stream
}

type testStreamBidirectionalStreamServer struct {
	server.Stream
}

func (x *testStreamBidirectionalStreamServer) Send(m *HelloRsp) error {
	return x.Stream.SendMsg(m)
}

func (x *testStreamBidirectionalStreamServer) Recv() (*HelloReq, error) {
	m := new(HelloReq)
	if err := x.Stream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TestStreamServer_ServiceDesc descriptor for server.RegisterService
var TestStreamServer_ServiceDesc = server.ServiceDesc{
	ServiceName:  "trpc.examples.stream.TestStream",
	HandlerType:  ((*TestStreamService)(nil)),
	StreamHandle: stream.NewStreamDispatcher(),
	Methods:      []server.Method{},
	Streams: []server.StreamDesc{
		{
			StreamName:    "/trpc.examples.stream.TestStream/ClientStream",
			Handler:       TestStreamService_ClientStream_Handler,
			ServerStreams: false,
		},
		{
			StreamName:    "/trpc.examples.stream.TestStream/ServerStream",
			Handler:       TestStreamService_ServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "/trpc.examples.stream.TestStream/BidirectionalStream",
			Handler:       TestStreamService_BidirectionalStream_Handler,
			ServerStreams: true,
		},
	},
}

// RegisterTestStreamService register service
func RegisterTestStreamService(s server.Service, svr TestStreamService) {
	if err := s.Register(&TestStreamServer_ServiceDesc, svr); err != nil {
		panic(fmt.Sprintf("TestStream register error:%v", err))
	}
}

// START --------------------------------- Default Unimplemented Server Service --------------------------------- START

type UnimplementedTestStream struct{}

// ClientStream Defined Client-side streaming RPC
//
//	Add stream in front of HelloReq
func (s *UnimplementedTestStream) ClientStream(stream TestStream_ClientStreamServer) error {
	return errors.New("rpc ClientStream of service TestStream is not implemented")
} // ServerStream Defined Server-side streaming RPC
// Add stream in front of HelloRsp
func (s *UnimplementedTestStream) ServerStream(req *HelloReq, stream TestStream_ServerStreamServer) error {
	return errors.New("rpc ServerStream of service TestStream is not implemented")
} // BidirectionalStream Bidirectional streaming RPC
// Add stream in front of HelloReq and HelloRsp
func (s *UnimplementedTestStream) BidirectionalStream(stream TestStream_BidirectionalStreamServer) error {
	return errors.New("rpc BidirectionalStream of service TestStream is not implemented")
}

// END --------------------------------- Default Unimplemented Server Service --------------------------------- END

// END ======================================= Server Service Definition ======================================= END

// START ======================================= Client Service Definition ======================================= START

// TestStreamClientProxy defines service client proxy
type TestStreamClientProxy interface {
	// ClientStream Defined Client-side streaming RPC
	//  Add stream in front of HelloReq
	ClientStream(ctx context.Context, opts ...client.Option) (TestStream_ClientStreamClient, error)
	// ServerStream Defined Server-side streaming RPC
	//  Add stream in front of HelloRsp
	ServerStream(ctx context.Context, req *HelloReq, opts ...client.Option) (TestStream_ServerStreamClient, error)
	// BidirectionalStream Bidirectional streaming RPC
	//  Add stream in front of HelloReq and HelloRsp
	BidirectionalStream(ctx context.Context, opts ...client.Option) (TestStream_BidirectionalStreamClient, error)
}

type TestStreamClientProxyImpl struct {
	client       client.Client
	streamClient stream.Client
	opts         []client.Option
}

var NewTestStreamClientProxy = func(opts ...client.Option) TestStreamClientProxy {
	return &TestStreamClientProxyImpl{client: client.DefaultClient, streamClient: stream.DefaultStreamClient, opts: opts}
}

func (c *TestStreamClientProxyImpl) ClientStream(ctx context.Context, opts ...client.Option) (TestStream_ClientStreamClient, error) {
	ctx, msg := codec.WithCloneMessage(ctx)

	msg.WithClientRPCName("/trpc.examples.stream.TestStream/ClientStream")
	msg.WithCalleeServiceName(TestStreamServer_ServiceDesc.ServiceName)
	msg.WithCalleeApp("examples")
	msg.WithCalleeServer("stream")
	msg.WithCalleeService("TestStream")
	msg.WithCalleeMethod("ClientStream")
	msg.WithSerializationType(codec.SerializationTypePB)

	clientStreamDesc := &client.ClientStreamDesc{}
	clientStreamDesc.StreamName = "/trpc.examples.stream.TestStream/ClientStream"
	clientStreamDesc.ClientStreams = true
	clientStreamDesc.ServerStreams = false

	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)

	stream, err := c.streamClient.NewStream(ctx, clientStreamDesc, "/trpc.examples.stream.TestStream/ClientStream", callopts...)
	if err != nil {
		return nil, err
	}
	x := &testStreamClientStreamClient{stream}
	return x, nil
}

type TestStream_ClientStreamClient interface {
	Send(*HelloReq) error
	CloseAndRecv() (*HelloRsp, error)
	client.ClientStream
}

type testStreamClientStreamClient struct {
	client.ClientStream
}

func (x *testStreamClientStreamClient) Send(m *HelloReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testStreamClientStreamClient) CloseAndRecv() (*HelloRsp, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(HelloRsp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *TestStreamClientProxyImpl) ServerStream(ctx context.Context, req *HelloReq, opts ...client.Option) (TestStream_ServerStreamClient, error) {
	ctx, msg := codec.WithCloneMessage(ctx)

	msg.WithClientRPCName("/trpc.examples.stream.TestStream/ServerStream")
	msg.WithCalleeServiceName(TestStreamServer_ServiceDesc.ServiceName)
	msg.WithCalleeApp("examples")
	msg.WithCalleeServer("stream")
	msg.WithCalleeService("TestStream")
	msg.WithCalleeMethod("ServerStream")
	msg.WithSerializationType(codec.SerializationTypePB)

	clientStreamDesc := &client.ClientStreamDesc{}
	clientStreamDesc.StreamName = "/trpc.examples.stream.TestStream/ServerStream"
	clientStreamDesc.ClientStreams = false
	clientStreamDesc.ServerStreams = true

	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)

	stream, err := c.streamClient.NewStream(ctx, clientStreamDesc, "/trpc.examples.stream.TestStream/ServerStream", callopts...)
	if err != nil {
		return nil, err
	}
	x := &testStreamServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(req); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TestStream_ServerStreamClient interface {
	Recv() (*HelloRsp, error)
	client.ClientStream
}

type testStreamServerStreamClient struct {
	client.ClientStream
}

func (x *testStreamServerStreamClient) Recv() (*HelloRsp, error) {
	m := new(HelloRsp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *TestStreamClientProxyImpl) BidirectionalStream(ctx context.Context, opts ...client.Option) (TestStream_BidirectionalStreamClient, error) {
	ctx, msg := codec.WithCloneMessage(ctx)

	msg.WithClientRPCName("/trpc.examples.stream.TestStream/BidirectionalStream")
	msg.WithCalleeServiceName(TestStreamServer_ServiceDesc.ServiceName)
	msg.WithCalleeApp("examples")
	msg.WithCalleeServer("stream")
	msg.WithCalleeService("TestStream")
	msg.WithCalleeMethod("BidirectionalStream")
	msg.WithSerializationType(codec.SerializationTypePB)

	clientStreamDesc := &client.ClientStreamDesc{}
	clientStreamDesc.StreamName = "/trpc.examples.stream.TestStream/BidirectionalStream"
	clientStreamDesc.ClientStreams = true
	clientStreamDesc.ServerStreams = true

	callopts := make([]client.Option, 0, len(c.opts)+len(opts))
	callopts = append(callopts, c.opts...)
	callopts = append(callopts, opts...)

	stream, err := c.streamClient.NewStream(ctx, clientStreamDesc, "/trpc.examples.stream.TestStream/BidirectionalStream", callopts...)
	if err != nil {
		return nil, err
	}
	x := &testStreamBidirectionalStreamClient{stream}
	return x, nil
}

type TestStream_BidirectionalStreamClient interface {
	Send(*HelloReq) error
	Recv() (*HelloRsp, error)
	client.ClientStream
}

type testStreamBidirectionalStreamClient struct {
	client.ClientStream
}

func (x *testStreamBidirectionalStreamClient) Send(m *HelloReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *testStreamBidirectionalStreamClient) Recv() (*HelloRsp, error) {
	m := new(HelloRsp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// END ======================================= Client Service Definition ======================================= END
