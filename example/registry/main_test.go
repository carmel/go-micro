package main

import (
	"context"
	"log"
	"testing"
	"time"

	ms "github.com/carmel/microservices"
	pb "github.com/carmel/microservices/example/testdata/helloworld"
	"github.com/carmel/microservices/midware/recovery"
	"github.com/carmel/microservices/registry/etcd"
	"github.com/carmel/microservices/transport/grpc"
	"github.com/carmel/microservices/transport/http"
	clientv3 "go.etcd.io/etcd/client/v3"
	srcgrpc "google.golang.org/grpc"
)

func TestServer(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	httpSrv := http.NewServer(
		http.Address(":8000"),
		http.Midware(
			recovery.Recovery(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		grpc.Midware(
			recovery.Recovery(),
		),
	)

	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	r := etcd.New(client)
	app := ms.New(
		ms.Name("helloworld"),
		ms.Server(
			httpSrv,
			grpcSrv,
		),
		ms.Registrar(r),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func TestClient(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(cli)

	connGRPC, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connGRPC.Close()

	connHTTP, err := http.NewClient(
		context.Background(),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithDiscovery(r),
		http.WithBlock(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer connHTTP.Close()

	for {
		callHTTP(connHTTP)
		callGRPC(connGRPC)
		time.Sleep(time.Second)
	}
}

func callGRPC(conn *srcgrpc.ClientConn) {
	client := pb.NewGreeterClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(conn *http.Client) {
	client := pb.NewGreeterHTTPClient(conn)
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %+v\n", reply)
}
