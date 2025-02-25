package config

import (
	"context"
	"flag"
	"log"
	"testing"
	"time"

	"go-micro/config"
	"go-micro/config/etcd"
	"go-micro/config/file"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "config.yaml", "config path, eg: -conf config.yaml")
}

func TestYamlAndEtcd(t *testing.T) {
	flag.Parse()

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = client.Close()
	}()

	testKey := "config.json"
	if _, err = client.Put(context.Background(), testKey, `{"service":{"name":"etcd server name"}}`); err != nil {
		t.Fatal(err)
	}

	source, err := etcd.New(client, etcd.WithPath(testKey))
	if err != nil {
		t.Fatal(err)
	}

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
			// 注意：WithSource加入的顺序影响配置读取覆盖的顺序
			source,
		),
	)

	if err := c.Load(); err != nil {
		panic(err)
	}

	// Defines the config JSON Field
	var v struct {
		Service struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"service"`
		Http struct {
			Server struct {
				Address string        `json:"address"`
				Timeout time.Duration `json:"timeout"`
			} `json:"server"`
		} `json:"http"`
	}

	// Unmarshal the config to struct
	if err := c.Scan(&v); err != nil {
		panic(err)
	}
	log.Printf("config: %+v", v)

	// Get a value associated with the key
	name, err := c.Value("service.name").String()
	if err != nil {
		panic(err)
	}
	log.Printf("service: %s", name)

	// watch key
	if err := c.Watch("service.name", func(key string, value config.Value) {
		log.Printf("config changed: %s = %v\n", key, value)
	}); err != nil {
		panic(err)
	}

	<-make(chan struct{})
}
