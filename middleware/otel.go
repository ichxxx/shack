package middleware

import (
	"fmt"

	"github.com/ichxxx/shack"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerKey  = "otel-go-contrib-tracer"
	tracerName = "otel-shack"
)

type config struct {
	TracerProvider oteltrace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

// OtelOption specifies instrumentation configuration options.
type OtelOption interface {
	apply(*config)
}

type optionFunc func(*config)

func (o optionFunc) apply(c *config) {
	o(c)
}

func OpenTelemetry(service string, opts ...OtelOption) shack.Handler {
	cfg := config{}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if cfg.TracerProvider == nil {
		cfg.TracerProvider = otel.GetTracerProvider()
	}
	tracer := cfg.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion(SemVersion()),
	)
	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	return func(ctx *shack.Context) {
		ctx.Set(tracerKey, tracer)
		savedCtx := ctx.Request.Context()
		defer func() {
			ctx.Request.Request = ctx.Request.WithContext(savedCtx)
		}()
		c := cfg.Propagators.Extract(savedCtx, propagation.HeaderCarrier(ctx.Request.Request.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", ctx.Request.Request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(ctx.Request.Request)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(service, ctx.Request.Path(), ctx.Request.Request)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		spanName := ctx.Request.Path()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", ctx.Request.Method())
		}
		c, span := tracer.Start(c, spanName, opts...)
		defer span.End()

		// pass the span through the request context
		ctx.Request.Request = ctx.Request.Request.WithContext(c)

		// serve the request to the next middleware
		ctx.Next()

		status := ctx.Response.StatusCode
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if ctx.Err != nil {
			span.SetAttributes(attribute.String("error", ctx.Err.Error()))
		}
	}
}

// WithPropagators specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) OtelOption {
	return optionFunc(func(cfg *config) {
		if propagators != nil {
			cfg.Propagators = propagators
		}
	})
}

// WithTracerProvider specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) OtelOption {
	return optionFunc(func(cfg *config) {
		if provider != nil {
			cfg.TracerProvider = provider
		}
	})
}

// Version is the current release version of the gin instrumentation.
func Version() string {
	return "0.29.0"
	// This string is updated by the pre_release.sh script during release
}

// SemVersion is the semantic version to be supplied to tracer/meter creation.
func SemVersion() string {
	return "semver:" + Version()
}
