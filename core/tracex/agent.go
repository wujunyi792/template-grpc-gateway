package tracex

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"pinnacle-primary-be/core/langx"
	"pinnacle-primary-be/core/logx"
	"sync"
)

const (
	kindOtlpGrpc = "otlpgrpc"
	kindOtlpHttp = "otlphttp"
)

var (
	agents = make(map[string]langx.PlaceholderType)
	lock   sync.Mutex
	tp     *sdktrace.TracerProvider
)

// StartAgent starts an opentelemetry agent.
func StartAgent(c Config) {
	lock.Lock()
	defer lock.Unlock()

	_, ok := agents[c.Endpoint]
	if ok {
		return
	}

	// if error happens, let later calls run.
	if err := startAgent(c); err != nil {
		return
	}

	agents[c.Endpoint] = langx.Placeholder
}

// StopAgent shuts down the span processors in the order they were registered.
func StopAgent() {
	_ = tp.Shutdown(context.Background())
}

func createExporter(c Config) (sdktrace.SpanExporter, error) {
	headers := map[string]string{"Authentication": c.Auth}
	// Just support jaeger and zipkin now, more for later
	switch c.Proto {
	case kindOtlpGrpc:
		// Always treat trace exporter as optional component, so we use nonblock here,
		// otherwise this would slow down app start up even set a dial timeout here when
		// endpoint can not reach.
		// If the connection not dial success, the global otel ErrorHandler will catch error
		// when reporting data like other exporters.
		return otlptracegrpc.New(
			context.Background(),
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithHeaders(headers),
			otlptracegrpc.WithEndpoint(c.Endpoint),
		)
	case kindOtlpHttp:
		// Not support flexible configuration now.
		return otlptracehttp.New(
			context.Background(),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(c.Endpoint),
		)
	default:
		return nil, fmt.Errorf("unknown exporter: %s", c.Proto)
	}
}

func startAgent(c Config) error {
	opts := []sdktrace.TracerProviderOption{
		// Set the sampling rate based on the parent span to 100%
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(c.Sampler))),
		// Record information about this application in a Resource.
		sdktrace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(c.Name))),
	}

	if len(c.Endpoint) > 0 {
		exp, err := createExporter(c)
		if err != nil {
			return err
		}

		// Always be sure to batch in production.
		opts = append(opts, sdktrace.WithBatcher(exp))
	}

	tp = sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		logx.Errorf("[otel] error: %v", err)
	}))

	return nil
}
