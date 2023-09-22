package main

import (
	"context"
	"fmt"
	"log"

	ms "github.com/carmel/microservices"
	pb "github.com/carmel/microservices/example/testdata/helloworld"
	"github.com/carmel/microservices/midware/recovery"
	"github.com/carmel/microservices/registry/etcd"
	"github.com/carmel/microservices/transport/grpc"
	"github.com/carmel/microservices/transport/http"
	etcdclient "go.etcd.io/etcd/client/v3"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}

func main() {
	client, err := etcdclient.New(etcdclient.Config{
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
