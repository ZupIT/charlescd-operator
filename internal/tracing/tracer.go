package tracing

import (
	"io"
	"runtime"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

// Initialize create an instance of Jaeger Tracer and sets it as GlobalTracer.
func Initialize(options ...TracerOptionFunc) (io.Closer, error) {
	cfg := new(config.Configuration)

	for _, fn := range options {
		if fn == nil {
			continue
		}
		if err := fn(cfg); err != nil {
			return nil, err
		}
	}

	cfg, err := cfg.FromEnv()
	if err != nil {
		return nil, err
	}

	// if the service name has not been defined, infer from the caller's package name
	if cfg.ServiceName == "" {
		pc, _, _, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()
		lastSlash := strings.LastIndexByte(funcName, '/')
		if lastSlash < 0 {
			lastSlash = 0
		} else {
			lastSlash++
		}
		firstDot := strings.IndexByte(funcName[lastSlash:], '.') + lastSlash
		cfg.ServiceName = funcName[lastSlash:firstDot]
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}
