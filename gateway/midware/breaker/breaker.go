package breaker

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	config "github.com/carmel/microservices/gateway/api/gateway/config/v1"
	v1 "github.com/carmel/microservices/gateway/api/gateway/midware/breaker/v1"
	"github.com/carmel/microservices/gateway/client"
	"github.com/carmel/microservices/gateway/midware"
	"github.com/carmel/microservices/gateway/proxy/condition"
	"github.com/carmel/microservices/logger"
	"github.com/carmel/microservices/pkg/breaker"
	"github.com/carmel/microservices/pkg/breaker/sre"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/rand"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func Init(clientFactory client.Factory) {
	breakerFactory := New(clientFactory)
	midware.RegisterV2("breaker", breakerFactory)
	prometheus.MustRegister(_metricDeniedTotal)
}

var (
	_metricDeniedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "go",
		Subsystem: "gateway",
		Name:      "requests_circuit_breaker_denied_total",
		Help:      "The total number of denied requests",
	}, []string{"protocol", "method", "path", "service", "basePath"})
)

type ratioTrigger struct {
	*v1.CircuitBreaker_Ratio
	lock sync.Mutex
	rand *rand.Rand
}

func newRatioTrigger(in *v1.CircuitBreaker_Ratio) *ratioTrigger {
	return &ratioTrigger{
		CircuitBreaker_Ratio: in,
		rand:                 rand.New(rand.NewSource(uint64(time.Now().UnixNano()))),
	}
}

func (r *ratioTrigger) Allow() error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.rand.Int63n(10000) < r.Ratio {
		return nil
	}
	return breaker.ErrNotAllowed
}
func (*ratioTrigger) MarkSuccess() {}
func (*ratioTrigger) MarkFailed()  {}

type nopTrigger struct{}

func (nopTrigger) Allow() error { return nil }
func (nopTrigger) MarkSuccess() {}
func (nopTrigger) MarkFailed()  {}

func makeBreakerTrigger(in *v1.CircuitBreaker) breaker.CircuitBreaker {
	switch trigger := in.Trigger.(type) {
	case *v1.CircuitBreaker_SuccessRatio:
		var opts []sre.Option
		if trigger.SuccessRatio.Bucket != 0 {
			opts = append(opts, sre.WithBucket(int(trigger.SuccessRatio.Bucket)))
		}
		if trigger.SuccessRatio.Request != 0 {
			opts = append(opts, sre.WithRequest(int64(trigger.SuccessRatio.Request)))
		}
		if trigger.SuccessRatio.Success != 0 {
			opts = append(opts, sre.WithSuccess(trigger.SuccessRatio.Success))
		}
		if trigger.SuccessRatio.Window != nil {
			opts = append(opts, sre.WithWindow(trigger.SuccessRatio.Window.AsDuration()))
		}
		return sre.NewBreaker(opts...)
	case *v1.CircuitBreaker_Ratio:
		return newRatioTrigger(trigger)
	default:
		logger.Warnf("Unrecoginzed circuit breaker trigger: %+v", trigger)
		return nopTrigger{}
	}
}

func makeOnBreakHandler(in *v1.CircuitBreaker, factory client.Factory) (http.RoundTripper, io.Closer, error) {
	switch action := in.Action.(type) {
	case *v1.CircuitBreaker_BackupService:
		logger.Infof("Making backup service as on break handler: %+v", action)
		client, err := factory(action.BackupService.Endpoint)
		if err != nil {
			return nil, nil, err
		}
		return client, client, nil
	case *v1.CircuitBreaker_ResponseData:
		logger.Infof("Making static response data as on break handler: %+v", action)
		return midware.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: int(action.ResponseData.StatusCode),
				Header:     http.Header{},
			}
			for _, h := range action.ResponseData.Header {
				resp.Header[h.Key] = h.Value
			}
			resp.Body = io.NopCloser(bytes.NewReader(action.ResponseData.Body))
			return resp, nil
		}), io.NopCloser(nil), nil
	default:
		logger.Warnf("Unrecoginzed circuit breaker aciton: %+v", action)
		return midware.RoundTripperFunc(func(*http.Request) (*http.Response, error) {
			// TBD: on break response
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Header:     http.Header{},
				Body:       io.NopCloser(&bytes.Buffer{}),
			}, nil
		}), io.NopCloser(nil), nil
	}
}

func isSuccessResponse(conditions []condition.Condition, resp *http.Response) bool {
	return condition.JudgeConditons(conditions, resp, true)
}

func deniedRequestIncr(req *http.Request) {
	labels, ok := midware.MetricsLabelsFromContext(req.Context())
	if ok {
		_metricDeniedTotal.WithLabelValues(labels.Protocol(), labels.Method(), labels.Path(), labels.Service(), labels.BasePath()).Inc()
		return
	}
}

func New(factory client.Factory) midware.FactoryV2 {
	return func(c *config.Midware) (midware.MidwareV2, error) {
		options := &v1.CircuitBreaker{}
		if c.Options != nil {
			if err := anypb.UnmarshalTo(c.Options, options, proto.UnmarshalOptions{Merge: true}); err != nil {
				return nil, err
			}
		}
		breaker := makeBreakerTrigger(options)
		onBreakHandler, closer, err := makeOnBreakHandler(options, factory)
		if err != nil {
			return nil, err
		}
		assertCondtions, err := condition.ParseConditon(options.AssertCondtions...)
		if err != nil {
			return nil, err
		}

		return midware.NewWithCloser(func(next http.RoundTripper) http.RoundTripper {
			return midware.RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
				if err := breaker.Allow(); err != nil {
					// rejected
					// NOTE: when client reject requests locally,
					// continue add counter let the drop ratio higher.
					breaker.MarkFailed()
					deniedRequestIncr(req)
					return onBreakHandler.RoundTrip(req)
				}
				resp, err := next.RoundTrip(req)
				if err != nil {
					breaker.MarkFailed()
					return nil, err
				}
				if !isSuccessResponse(assertCondtions, resp) {
					breaker.MarkFailed()
					return resp, nil
				}
				breaker.MarkSuccess()
				return resp, nil
			})
		}, closer), nil
	}
}
