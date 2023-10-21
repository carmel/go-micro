package main

import (
	"context"
	"log"
	"testing"
	"time"

	pb "go-micro/example/testdata/helloworld"
	"go-micro/logger"
	"go-micro/midware/recovery"
	"go-micro/registry/etcd"
	"go-micro/selector/filter"
	"go-micro/transport/grpc"
	"go-micro/transport/http"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestServer(t *testing.T) {

	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}

	slog, _ := logger.NewSlogger(logger.Options{LogPath: "log/ms.log"})

	go runServer("1.0", slog, client, 8000)
	go runServer("1.0", slog, client, 8010)

	runServer("2.0", slog, client, 8020)
}

func TestClient(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(cli)

	// selector.SetGlobalSelector(wrr.NewBuilder())

	// new grpc client
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///helloworld"),
		grpc.WithDiscovery(r),
		// 由于gRPC框架的限制只能使用全局balancer+filter的方式来实现selector
		// 这里使用weighted round robin算法的balancer+静态version=1.0.0的Filter
		grpc.WithNodeFilter(
			filter.Version("1.0"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	gClient := pb.NewGreeterClient(conn)

	// new http client
	hConn, err := http.NewClient(
		context.Background(),
		http.WithMidware(
			recovery.Recovery(),
		),
		http.WithEndpoint("discovery:///helloworld"),
		http.WithDiscovery(r),
		// 这里使用p2c算法的balancer+静态version=2.0.0的Filter组成一个selector
		http.WithNodeFilter(filter.Version("2.0")),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer hConn.Close()
	hClient := pb.NewGreeterHTTPClient(hConn)

	for {
		time.Sleep(time.Second)
		callGRPC(gClient)
		callHTTP(hClient)
	}
}

func callGRPC(client pb.GreeterClient) {
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[grpc] SayHello %+v\n", reply)
}

func callHTTP(client pb.GreeterHTTPClient) {
	reply, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "kratos"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[http] SayHello %s\n", reply.Message)
}
