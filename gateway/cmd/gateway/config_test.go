package main

import (
	"testing"

	"github.com/carmel/microservices/config"
	"github.com/carmel/microservices/config/file"
	configv1 "github.com/carmel/microservices/gateway/api/config/v1"
	circuitbreakerv1 "github.com/carmel/microservices/gateway/api/midware/breaker/v1"
	corsv1 "github.com/carmel/microservices/gateway/api/midware/cors/v1"
	rewritev1 "github.com/carmel/microservices/gateway/api/midware/rewrite/v1"
	tracingv1 "github.com/carmel/microservices/gateway/api/midware/tracing/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

func equalTo() *configv1.Gateway {
	return &configv1.Gateway{
		Name:    "helloworld",
		Version: "v1",
		// Hosts: []string{
		// 	"localhost",
		// 	"127.0.0.1",
		// },
		Endpoints: []*configv1.Endpoint{
			{
				Path:     "/helloworld/*",
				Protocol: configv1.Protocol_HTTP,
				Host:     "localhost",
				Timeout:  &durationpb.Duration{Seconds: 1},
				Backends: []*configv1.Backend{
					{
						Target: "127.0.0.1:8000",
					},
				},
				Middlewares: []*configv1.Midware{
					{
						Name: "breaker",
						Options: asAny(&circuitbreakerv1.CircuitBreaker{
							Trigger: &circuitbreakerv1.CircuitBreaker_SuccessRatio{
								SuccessRatio: &circuitbreakerv1.SuccessRatio{
									Success: 0.6,
									Request: 1,
									Bucket:  10,
									Window:  &durationpb.Duration{Seconds: 3},
								},
							},
							Action: &circuitbreakerv1.CircuitBreaker_BackupService{
								BackupService: &circuitbreakerv1.BackupService{
									Endpoint: &configv1.Endpoint{
										Backends: []*configv1.Backend{
											{
												Target: "127.0.0.1:8001",
											},
										},
									},
								},
							},
							AssertCondtions: []*configv1.Condition{
								{
									Condition: &configv1.Condition_ByStatusCode{
										ByStatusCode: "200",
									},
								},
							},
						}),
					},
					{
						Name:    "rewrite",
						Options: asAny(&rewritev1.Rewrite{}),
					},
				},
			},
			{
				Path:     "/helloworld.Greeter/*",
				Method:   "POST",
				Protocol: configv1.Protocol_GRPC,
				Timeout:  &durationpb.Duration{Seconds: 1},
				Backends: []*configv1.Backend{
					{
						Target: "127.0.0.1:9000",
					},
				},
				Retry: &configv1.Retry{
					Attempts:      3,
					PerTryTimeout: &durationpb.Duration{Nanos: 100000000},
					Conditions: []*configv1.Condition{
						{Condition: &configv1.Condition_ByStatusCode{ByStatusCode: "502-504"}},
						{Condition: &configv1.Condition_ByHeader{ByHeader: &configv1.ConditionHeader{
							Name:  "Grpc-Status",
							Value: "14",
						}}},
					},
				},
			},
		},
		Middlewares: []*configv1.Midware{
			{
				Name: "tracing",
				Options: asAny(&tracingv1.Tracing{
					HttpEndpoint: "localhost:4318",
				}),
			},
			{
				Name: "logging",
			},
			{
				Name: "transcoder",
			},
			{
				Name: "cors",
				Options: asAny(&corsv1.Cors{
					AllowCredentials: true,
					AllowOrigins:     []string{".google.com"},
					AllowMethods:     []string{"GET", "POST", "OPTIONS"},
				}),
			},
		},
	}
}

func asAny(in proto.Message) *anypb.Any {
	out, err := anypb.New(in)
	if err != nil {
		panic(err)
	}
	return out
}

func TestConfigUnmarshaler(t *testing.T) {
	cfg := config.New(
		config.WithSource(
			file.NewSource("config.yaml"),
		),
	)
	if err := cfg.Load(); err != nil {
		t.Fatal(err)
	}
	gateway := &configv1.Gateway{}
	if err := cfg.Scan(gateway); err != nil {
		t.Fatal(err)
	}

	left, err := protojson.Marshal(gateway)
	if err != nil {
		t.Fatal(err)
	}
	right, err := protojson.Marshal(equalTo())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("gateway config:\nloaded: %s\nshould equal to: %s\n", left, right)

	if !proto.Equal(gateway, equalTo()) {
		t.Errorf("inconsistent gateway config")
	}
}

// func TestLog(t *testing.T) {
// 	log := logger.With("test", "gateway")
// 	log.Error("some error %s", "sdfds", "22")
// }
