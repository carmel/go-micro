package main

import (
	"context"
	"log"
	"testing"

	ms "github.com/carmel/microservices"
	pb "github.com/carmel/microservices/example/testdata/helloworld"
	"github.com/carmel/microservices/metadata"
	mmd "github.com/carmel/microservices/midware/metadata"
	"github.com/carmel/microservices/transport/grpc"
	"github.com/carmel/microservices/transport/http"
)

func TestServer(t *testing.T) {
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Midware(
			mmd.Server(),
		))
	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Midware(
			mmd.Server(),
		),
	)

	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	app := ms.New(
		ms.Name("helloworld"),
		ms.Server(
			httpSrv,
			grpcSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func TestClient(t *testing.T) {
	callHTTP()
	callGRPC()
}

func callHTTP() {
	conn, err := http.NewClient(
		context.Background(),
		http.WithMidware(
			mmd.Client(),
		),
		http.WithEndpoint("127.0.0.1:8000"),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterHTTPClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToClientContext(ctx, "x-md-global-extra", "2233")
	reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply)
}

func callGRPC() {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMidware(
			mmd.Client(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToClientContext(ctx, "x-md-global-extra", "2233")
	reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v \n", reply)
}
