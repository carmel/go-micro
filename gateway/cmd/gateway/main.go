package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	ms "github.com/carmel/go-micro"
	"github.com/carmel/go-micro/gateway/client"
	"github.com/carmel/go-micro/gateway/config"
	configLoader "github.com/carmel/go-micro/gateway/config/config-loader"
	"github.com/carmel/go-micro/gateway/discovery"
	"github.com/carmel/go-micro/gateway/midware"
	"github.com/carmel/go-micro/gateway/proxy"
	"github.com/carmel/go-micro/gateway/proxy/debug"
	"github.com/carmel/go-micro/gateway/server"

	_ "net/http/pprof"

	_ "github.com/carmel/go-micro/gateway/discovery/consul"
	_ "github.com/carmel/go-micro/gateway/midware/bbr"
	"github.com/carmel/go-micro/gateway/midware/breaker"
	_ "github.com/carmel/go-micro/gateway/midware/cors"
	_ "github.com/carmel/go-micro/gateway/midware/logging"
	_ "github.com/carmel/go-micro/gateway/midware/rewrite"
	_ "github.com/carmel/go-micro/gateway/midware/tracing"
	_ "github.com/carmel/go-micro/gateway/midware/transcoder"
	_ "go.uber.org/automaxprocs"

	"github.com/carmel/go-micro/logger"
	"github.com/carmel/go-micro/registry"
	"github.com/carmel/go-micro/transport"
	"golang.org/x/exp/rand"
)

var (
	ctrlName     string
	ctrlService  string
	discoveryDSN string
	proxyAddrs   = newSliceVar(":8080")
	proxyConfig  string
	withDebug    bool
)

type sliceVar struct {
	val        []string
	defaultVal []string
}

func newSliceVar(defaultVal ...string) sliceVar {
	return sliceVar{defaultVal: defaultVal}
}
func (s *sliceVar) Get() []string {
	if len(s.val) <= 0 {
		return s.defaultVal
	}
	return s.val
}
func (s *sliceVar) Set(val string) error {
	s.val = append(s.val, val)
	return nil
}
func (s *sliceVar) String() string { return fmt.Sprintf("%+v", *s) }

func init() {
	rand.Seed(uint64(time.Now().Nanosecond()))

	flag.BoolVar(&withDebug, "debug", false, "enable debug handlers")
	flag.Var(&proxyAddrs, "addr", "proxy address, eg: -addr 0.0.0.0:8080")
	flag.StringVar(&proxyConfig, "conf", "config.yaml", "config path, eg: -conf config.yaml")
	flag.StringVar(&ctrlName, "ctrl.name", os.Getenv("ADVERTISE_NAME"), "control gateway name, eg: gateway")
	flag.StringVar(&ctrlService, "ctrl.service", "", "control service host, eg: http://127.0.0.1:8000")
	flag.StringVar(&discoveryDSN, "discovery.dsn", "", "discovery dsn, eg: consul://127.0.0.1:7070?token=secret&datacenter=prod")
}

func makeDiscovery() registry.Discovery {
	if discoveryDSN == "" {
		return nil
	}
	d, err := discovery.Create(discoveryDSN)
	if err != nil {
		logger.Errorf("failed to create discovery: %v", err)
		return nil
	}
	return d
}

func main() {
	flag.Parse()

	clientFactory := client.NewFactory(makeDiscovery())
	p, err := proxy.New(clientFactory, midware.Create)
	if err != nil {
		logger.Errorf("failed to new proxy: %v", err)
		return
	}
	breaker.Init(clientFactory)

	ctx := context.Background()
	var ctrlLoader *configLoader.CtrlConfigLoader
	if ctrlService != "" {
		logger.Infof("setup control service to: %q", ctrlService)
		ctrlLoader = configLoader.New(ctrlName, ctrlService, proxyConfig)
		if err := ctrlLoader.Load(ctx); err != nil {
			logger.Errorf("failed to do initial load from control service: %v, using local config instead", err)
		}
		if err := ctrlLoader.LoadFeatures(ctx); err != nil {
			logger.Errorf("failed to do initial feature load from control service: %v, using default value instead", err)
		}
		go ctrlLoader.Run(ctx)
	}

	confLoader, err := config.NewFileLoader(proxyConfig)
	if err != nil {
		logger.Errorf("failed to create config file loader: %v", err)
		return
	}
	defer confLoader.Close()
	bc, err := confLoader.Load(context.Background())
	if err != nil {
		logger.Errorf("failed to load config: %v", err)
		return
	}

	if err := p.Update(bc); err != nil {
		logger.Errorf("failed to update service config: %v", err)
		return
	}
	reloader := func() error {
		bc, err := confLoader.Load(context.Background())
		if err != nil {
			logger.Errorf("failed to load config: %v", err)
			return err
		}
		if err := p.Update(bc); err != nil {
			logger.Errorf("failed to update service config: %v", err)
			return err
		}
		logger.Infof("config reloaded")
		return nil
	}
	confLoader.Watch(reloader)

	var serverHandler http.Handler = p
	if withDebug {
		debug.Register("proxy", p)
		debug.Register("config", confLoader)
		if ctrlLoader != nil {
			debug.Register("ctrl", ctrlLoader)
		}
		serverHandler = debug.MashupWithDebugHandler(p)
	}
	servers := make([]transport.Server, 0, len(proxyAddrs.Get()))
	for _, addr := range proxyAddrs.Get() {
		servers = append(servers, server.NewProxy(serverHandler, addr))
	}
	app := ms.New(
		ms.Name(bc.Name),
		ms.Context(ctx),
		ms.Server(
			servers...,
		),
	)
	if err := app.Run(); err != nil {
		logger.Errorf("failed to run servers: %v", err)
	}
}
