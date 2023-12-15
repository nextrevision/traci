package tracing

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"os"
	"strings"
)

const TraceParentKey = "TRACEPARENT"

func NewTraceProvider(ctx context.Context, serviceName string, resourceAttributes []attribute.KeyValue) *sdktrace.TracerProvider {
	exporter, err := newExporter(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR could not create span exporter: %v\n", err)
	}

	resources, err := resource.New(ctx,
		resource.WithAttributes(resourceAttributes...),
		resource.WithAttributes(semconv.ServiceName(serviceName)),
		resource.WithContainer(),
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithOS(),
		resource.WithProcessOwner(),
		resource.WithProcessPID(),
		resource.WithProcessRuntimeName(),
		resource.WithProcessRuntimeVersion(),
		resource.WithProcessRuntimeDescription(),
	)

	// Create provider using the exporter
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSyncer(exporter),
		sdktrace.WithResource(resources),
	)
}

func NewTracer(name string, provider trace.TracerProvider) trace.Tracer {
	otel.SetTracerProvider(provider)
	return otel.Tracer(name)
}

func newConsoleExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func newGrpcExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(ctx)
}

func newHttpExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracehttp.New(ctx)
}

// newExporter generates a new instance of sdktrace.SpanExporter based on the provided context.
// It reads the OTEL_EXPORTER_OTLP_ENDPOINT and OTEL_EXPORTER_OTLP_PROTOCOL environment variables to determine the exporter type.
// If OTEL_EXPORTER_OTLP_PROTOCOL is not set, it checks the endpoint to infer the protocol (grpc or http/json).
// If the protocol is grpc, it returns a new instance of grpc exporter created by newGrpcExporter.
// If the protocol is http, it returns a new instance of HTTP exporter created by newHttpExporter.
// If the protocol is console, it returns a new instance of console exporter created by newConsoleExporter.
// The context is passed to the selected exporter function for proper initialization and configuration.
func newExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	proto := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")

	if proto == "" {
		if strings.Contains(endpoint, ":4317") {
			proto = "grpc"
		} else if strings.Contains(endpoint, ":4318") {
			proto = "http/json"
		}
	}

	if strings.Contains(proto, "grpc") {
		return newGrpcExporter(ctx)
	} else if strings.Contains(proto, "http") {
		return newHttpExporter(ctx)
	} else if strings.Contains(proto, "console") {
		return newConsoleExporter()
	}

	// Return a no-op exporter to ensure the tracer does not panic and the command executes
	return &tracetest.NoopExporter{}, errors.New("could not determine OTLP protocol; set with env var OTEL_EXPORTER_OTLP_PROTOCOL")
}

// GenTraceParentString formats the TraceID and SpanID from the provided SpanContext and returns a W3C TraceParent string
func GenTraceParentString(spanContext trace.SpanContext) string {
	return fmt.Sprintf("00-%s-%s-00", spanContext.TraceID(), spanContext.SpanID())
}

// NewContextFromEnvTraceParent generates a new context with a trace parent extracted from the `TRACEPARENT` environment variable.
// If the `TRACEPARENT` environment variable is not present, the function returns the provided context.
//
// Example usage:
//
//	 os.Setenv("TRACEPARENT", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00")
//		ctx := NewContextFromEnvTraceParent(context.Background())
func NewContextFromEnvTraceParent(ctx context.Context) (context.Context, error) {
	val, present := os.LookupEnv(TraceParentKey)
	if !present {
		return ctx, fmt.Errorf("%s variable not found in environment", TraceParentKey)
	}

	carrier := propagation.MapCarrier{}
	carrier.Set("traceparent", val)

	tc := propagation.TraceContext{}
	newCtx := tc.Extract(ctx, carrier)
	if sc := trace.SpanContextFromContext(newCtx); !sc.IsValid() {
		return ctx, fmt.Errorf("%s variable is invalid", TraceParentKey)
	}

	return newCtx, nil
}

// NewContextFromDeterministicString generates a new context with a deterministic trace ID and span ID from the provided strings.
// The returned context will include the trace ID, span ID, and trace flags set.
//
// Example usage:
//
//	ctx := NewContextFromDeterministicString("foo", "bar")
//	tracer := otel.Tracer("test-tracer")
//	ctxSpan, span := tracer.Start(ctx, "test-span")
func NewContextFromDeterministicString(traceIDString string) context.Context {
	traceID, err := genTraceIDFromString(traceIDString)
	if err != nil {
		log.Fatalf("could not generate trace ID from string %s\n", traceIDString)
	}

	bytes := make([]byte, 8)
	rand.Read(bytes)
	spanID, err := genSpanIDFromString(hex.EncodeToString(bytes))
	if err != nil {
		log.Fatalf("could not generate trace ID from string %s\n", traceIDString)
	}

	return trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: 0x0,
	}))
}

// genTraceIDFromString generates an idempotent trace.TraceID from an arbitrary string.
func genTraceIDFromString(input string) (trace.TraceID, error) {
	var traceID [16]byte
	traceID = md5.Sum([]byte(input))
	return trace.TraceIDFromHex(fmt.Sprintf("%x", traceID)[:32])
}

// genSpanIDFromString generates an idempotent trace.SpanID from an arbitrary string.
func genSpanIDFromString(input string) (trace.SpanID, error) {
	var spanID [16]byte
	spanID = md5.Sum([]byte(input))
	return trace.SpanIDFromHex(fmt.Sprintf("%x", spanID)[:16])
}

// AttributeMapToKeyValue converts a map of attributes into a slice of attribute.KeyValue.
func AttributeMapToKeyValue(attributes map[string]string) []attribute.KeyValue {
	var spanAttributes []attribute.KeyValue
	for k, v := range attributes {
		spanAttributes = append(spanAttributes, attribute.String(k, v))
	}
	return spanAttributes
}
