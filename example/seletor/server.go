package main

import (
	"context"
	"fmt"
	"log"

	ms "go-micro"
	pb "go-micro/example/testdata/helloworld"
	"go-micro/logger"
	"go-micro/midware/logging"
	"go-micro/midware/recovery"
	"go-micro/registry/etcd"
	"go-micro/transport/grpc"
	"go-micro/transport/http"

	etcdclient "go.etcd.io/etcd/client/v3"
)

// server is used to implement pb.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements pb.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Welcome %+v!", in.Name)}, nil
}

func main() {

	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		log.Fatalf("[selector] new etcd client failed: %s", err)
	}

	var slog logger.Logger
	slog, err = logger.NewSlogger(logger.Options{LogPath: "log/ms.log"})
	if err != nil {
		log.Fatalf("[selector] new slogger failed: %s", err)
	}

	go runServer("1.0", slog, client, 8000)
	go runServer("1.0", slog, client, 8010)

	runServer("2.0", slog, client, 8020)
}

func runServer(version string, logger logger.Logger, client *etcdclient.Client, port int) {
	logger = logger.With(logger, "version", version, "port:", port)

	httpSrv := http.NewServer(
		http.Address(fmt.Sprintf(":%d", port)),
		http.Midware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)
	grpcSrv := grpc.NewServer(
		grpc.Address(fmt.Sprintf(":%d", port+1000)),
		grpc.Midware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	)

	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	r := etcd.New(client)
	app := ms.New(
		ms.Name("helloworld"),
		ms.Server(
			grpcSrv,
			httpSrv,
		),
		ms.Version(version),
		ms.Registrar(r),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
