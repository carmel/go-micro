package main

import (
	"context"

	"go-micro/client"
	"go-micro/examples/helloworld/pb"
	"go-micro/log"
)

func main() {
	c := pb.NewGreeterClientProxy(client.WithTarget("ip://127.0.0.1:8000"))
	rsp, err := c.Hello(context.Background(), &pb.HelloRequest{Msg: "world"})
	if err != nil {
		log.Error(err)
	}
	log.Info(rsp.Msg)
}
