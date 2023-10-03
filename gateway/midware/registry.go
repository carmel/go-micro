package midware

import (
	"errors"
	"strings"

	configv1 "github.com/carmel/microservices/gateway/api/gateway/config/v1"
	"github.com/carmel/microservices/logger"
	"github.com/prometheus/client_golang/prometheus"
)

var LOG logger.Logger //= logger.NewHelper(logger.With(logger.GetLogger(), "source", "midware"))
var globalRegistry = NewRegistry()
var _failedMiddlewareCreate = prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: "go",
	Subsystem: "gateway",
	Name:      "failed_middleware_create",
	Help:      "The total number of failed midware create",
}, []string{"name", "required"})

func Init(log logger.Logger) {
	prometheus.MustRegister(_failedMiddlewareCreate)
	LOG = log.With("source", "midware")
}

// ErrNotFound is midware not found.
var ErrNotFound = errors.New("Midware has not been registered")

// Registry is the interface for callers to get registered midware.
type Registry interface {
	Register(name string, factory Factory)
	RegisterV2(name string, factory FactoryV2)
	Create(cfg *configv1.Midware) (MidwareV2, error)
}

type middlewareRegistry struct {
	midware map[string]FactoryV2
}

// NewRegistry returns a new midware registry.
func NewRegistry() Registry {
	return &middlewareRegistry{
		midware: map[string]FactoryV2{},
	}
}

// Register registers one midware.
func (p *middlewareRegistry) Register(name string, factory Factory) {
	p.midware[createFullName(name)] = wrapFactory(factory)
}

func (p *middlewareRegistry) RegisterV2(name string, factory FactoryV2) {
	p.midware[createFullName(name)] = factory
}

// Create instantiates a midware based on `cfg`.
func (p *middlewareRegistry) Create(cfg *configv1.Midware) (MidwareV2, error) {
	if method, ok := p.getMiddleware(createFullName(cfg.Name)); ok {
		if cfg.Required {
			// If the midware is required, it must be created successfully.
			instance, err := method(cfg)
			if err != nil {
				_failedMiddlewareCreate.WithLabelValues(cfg.Name, "true").Inc()
				LOG.With("reason", "create_required_middleware_failed", "name", cfg.Name, "error", err, "config", cfg).Log(logger.ERROR, "Failed to create required midware")
				return nil, err
			}
			return instance, nil
		}
		instance, err := method(cfg)
		if err != nil {
			_failedMiddlewareCreate.WithLabelValues(cfg.Name, "false").Inc()
			LOG.With("reason", "create_optional_middleware_failed", "name", cfg.Name, "error", err, "config", cfg).Log(logger.ERROR, "Failed to create optional midware")
			return EmptyMiddleware, nil
		}
		return instance, nil
	}
	return nil, ErrNotFound
}

func (p *middlewareRegistry) getMiddleware(name string) (FactoryV2, bool) {
	nameLower := strings.ToLower(name)
	middlewareFn, ok := p.midware[nameLower]
	if ok {
		return middlewareFn, true
	}
	return nil, false
}

func createFullName(name string) string {
	return strings.ToLower("gateway.midware." + name)
}

// Register registers one midware.
func Register(name string, factory Factory) {
	globalRegistry.Register(name, factory)
}

// RegisterV2 registers one v2 midware.
func RegisterV2(name string, factory FactoryV2) {
	globalRegistry.RegisterV2(name, factory)
}

// Create instantiates a midware based on `cfg`.
func Create(cfg *configv1.Midware) (MidwareV2, error) {
	return globalRegistry.Create(cfg)
}
