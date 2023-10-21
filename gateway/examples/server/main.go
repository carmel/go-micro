package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	ms "go-micro"
	pb "go-micro/example/testdata/helloworld"
	"go-micro/midware/recovery"
	"go-micro/transport/grpc"
	"go-micro/transport/http"
)

var (
	httpAddr string
	grpcAddr string
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in.Name == "error" {
		return nil, context.DeadlineExceeded
	}
	return &pb.HelloReply{Message: in.GetName()}, nil
}

func init() {
	flag.StringVar(&httpAddr, "http.addr", ":8000", "server address, eg: 127.0.0.1:8000")
	flag.StringVar(&grpcAddr, "grpc.addr", ":9000", "server address, eg: 127.0.0.1:9000")
}

func main() {
	flag.Parse()
	s := &server{}
	httpSrv := http.NewServer(
		http.Address(httpAddr),
		http.Midware(
			recovery.Recovery(),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(grpcAddr),
		grpc.Midware(
			recovery.Recovery(),
		),
	)
	httpSrv.HandleFunc("/helloworld/header", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			fmt.Fprintf(w, "%s: %s\n", k, v)
		}
	})
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)
	app := ms.New(
		ms.Name("Helloworld"),
		ms.Server(
			httpSrv,
			grpcSrv,
		),
	)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
