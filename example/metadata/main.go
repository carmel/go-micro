package main

import (
	"context"
	"fmt"
	"log"

	ms "github.com/carmel/go-micro"
	pb "github.com/carmel/go-micro/example/testdata/helloworld"
	"github.com/carmel/go-micro/metadata"
	mmd "github.com/carmel/go-micro/midware/metadata"
	"github.com/carmel/go-micro/transport/grpc"
	"github.com/carmel/go-micro/transport/http"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
// Name is the name of the compiled software.
// Name = "helloworld"
// Version is the version of the compiled software.
// Version = "v1.0.0"
)

// server is used to implement pb.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements pb.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	var extra string
	if md, ok := metadata.FromServerContext(ctx); ok {
		extra = md.Get("x-md-global-extra")
	}
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %s extra_meta: %s", in.Name, extra)}, nil
}

func main() {

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
