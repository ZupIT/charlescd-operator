package tracing

import (
	"fmt"
	"net/url"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

type SamplerType string

const (
	ConstantSampler      SamplerType = jaeger.SamplerTypeConst
	ProbabilisticSampler SamplerType = jaeger.SamplerTypeProbabilistic
	RateLimitingSampler  SamplerType = jaeger.SamplerTypeRateLimiting
	RemoteSampler        SamplerType = jaeger.SamplerTypeRemote
)

type TracerOptionFunc func(*config.Configuration) error

func WithServiceName(service string) TracerOptionFunc {
	return func(o *config.Configuration) error {
		o.ServiceName = service
		return nil
	}
}

func WithSamplerType(samplerType SamplerType) TracerOptionFunc {
	return func(o *config.Configuration) error {
		if o.Sampler == nil {
			o.Sampler = new(config.SamplerConfig)
		}
		o.Sampler.Type = string(samplerType)
		return nil
	}
}

func WithSamplerParam(samplerParam float64) TracerOptionFunc {
	return func(o *config.Configuration) error {
		if o.Sampler == nil {
			o.Sampler = new(config.SamplerConfig)
		}
		o.Sampler.Param = samplerParam
		return nil
	}
}

func WithEndpoint(endpoint string) TracerOptionFunc {
	return func(o *config.Configuration) error {
		u, err := url.ParseRequestURI(endpoint)
		if err != nil {
			return fmt.Errorf("cannot parse collector endpoint %s: %w", endpoint, err)
		}
		if o.Reporter == nil {
			o.Reporter = new(config.ReporterConfig)
		}
		o.Reporter.CollectorEndpoint = u.String()
		return nil
	}
}
