package main

import (
	"context"

	trpc "go-micro"
	"go-micro/examples/helloworld/pb"
	"go-micro/log"
)

func main() {
	s := trpc.NewServer()
	pb.RegisterGreeterService(s, &Greeter{})
	if err := s.Serve(); err != nil {
		log.Error(err)
	}
}

type Greeter struct{}

func (g Greeter) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Infof("got hello request: %s", req.Msg)
	return &pb.HelloReply{Msg: "Hello " + req.Msg + "!"}, nil
}
