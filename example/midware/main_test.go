package main

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	ms "github.com/carmel/microservices"
	"github.com/carmel/microservices/logger"
	prom "github.com/carmel/microservices/metrics/prometheus"
	"github.com/carmel/microservices/midware/breaker"
	"github.com/carmel/microservices/midware/logging"
	"github.com/carmel/microservices/midware/metrics"
	"github.com/carmel/microservices/midware/ratelimit"
	"github.com/carmel/microservices/midware/recovery"
	"github.com/carmel/microservices/registry/etcd"
	"github.com/carmel/microservices/transport/grpc"
	"github.com/carmel/microservices/transport/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	clientv3 "go.etcd.io/etcd/client/v3"

	pb "github.com/carmel/microservices/example/testdata/helloworld"
	srcgrpc "google.golang.org/grpc"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = "metrics"
	// Version is the version of the compiled software.
	// Version = "v1.0.0"

	_metricSeconds = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_sec",
		Help:      "server requests duration(sec).",
		Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1},
	}, []string{"kind", "operation"})

	_metricRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "client",
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "The total number of processed requests",
	}, []string{"kind", "operation", "code", "reason"})
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements pb.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %+v", in.Name)}, nil
}

func init() {
	prometheus.MustRegister(_metricSeconds, _metricRequests)
}

func TestServer(t *testing.T) {
	slog := logger.NewSlogger(logger.Options{LogPath: "log/ms.log"})

	grpcSrv := grpc.NewServer(
		grpc.Address(":9000"),
		// server端中间件
		grpc.Midware(
			// 异常恢复
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					// do someting
					return nil
				}),
			),
			// 日志
			logging.Client(slog),
			// 限流器
			ratelimit.Server(),
			// 流量监控
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
			),
		),
	)

	httpSrv := http.NewServer(
		http.Address(":8000"),
		// server端中间件
		http.Midware(
			// 异常恢复
			recovery.Recovery(
				recovery.WithHandler(func(ctx context.Context, req, err interface{}) error {
					// do someting
					return nil
				}),
			),
			// 日志
			logging.Client(slog),
			// 限流器
			ratelimit.Server(),
			// 流量监控
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
			),
		),
	)
	httpSrv.Handle("/metrics", promhttp.Handler())

	s := &server{}
	pb.RegisterGreeterServer(grpcSrv, s)
	pb.RegisterGreeterHTTPServer(httpSrv, s)

	app := ms.New(
		ms.Name(Name),
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
		// client端中间件
		grpc.WithMidware(
			// 熔断器
			breaker.Client(),
		),
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
		// client端中间件
		http.WithMidware(
			// 熔断器
			breaker.Client(),
		),
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
