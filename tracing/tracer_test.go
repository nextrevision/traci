package tracing

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"
	"reflect"
	"testing"
)

func TestNewExporter(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name    string
		setVars func()
		wantErr bool
	}{
		{
			name: "Case for gRPC protocol (port 4317)",
			setVars: func() {
				os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317")
				os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "")
			},
			wantErr: false,
		},
		{
			name: "Case for HTTP protocol (port 4318)",
			setVars: func() {
				os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318")
				os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "")
			},
			wantErr: false,
		},
		{
			name: "Case for Console protocol",
			setVars: func() {
				os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
				os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "console")
			},
			wantErr: false,
		},
		{
			name: "Case for grpc protocol",
			setVars: func() {
				os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
				os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "grpc")
			},
			wantErr: false,
		},
		{
			name: "Case for http protocol",
			setVars: func() {
				os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
				os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "http/json")
			},
			wantErr: false,
		},
		{
			name: "Default noop exporter with err",
			setVars: func() {
				os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "")
				os.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setVars()
			got, err := newExporter(ctx)
			if (err != nil) != tc.wantErr {
				t.Errorf("error %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err != nil && tc.wantErr {
				return
			}
			if got == nil {
				t.Errorf("expected span exporter, got nil")
			}
		})
	}
}

func TestGenTraceParentString(t *testing.T) {
	tests := []struct {
		name        string
		spanContext trace.SpanContext
		want        string
	}{

		{
			name: "Valid input",
			spanContext: trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: mustTraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736"),
				SpanID:  mustSpanIDFromHex("00f067aa0ba902b7"),
			}),
			want: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
		},

		{
			name:        "Empty input",
			spanContext: trace.SpanContext{},
			want:        "00-00000000000000000000000000000000-0000000000000000-00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenTraceParentString(tt.spanContext)
			assert.Equal(t, tt.want, result)
		})
	}

}

func TestNewContextFromEnvTraceParent(t *testing.T) {
	// preparing a table driven test
	testCases := []struct {
		name        string // name of the test
		envKey      string
		envValue    string // 'traceparent' environment variable's value
		wantTraceID trace.TraceID
		wantSpanID  trace.SpanID
		wantErr     bool // flag indicating whether error is expected or not
	}{
		{
			name:        "No env var present",
			envKey:      "foo",
			envValue:    "",
			wantTraceID: trace.TraceID{},
			wantSpanID:  trace.SpanID{},
			wantErr:     true,
		},
		{
			name:        "Env var present",
			envKey:      TraceParentKey,
			envValue:    "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00",
			wantTraceID: mustTraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736"),
			wantSpanID:  mustSpanIDFromHex("00f067aa0ba902b7"),
			wantErr:     false,
		},
		{
			name:        "Env var invalid",
			envKey:      TraceParentKey,
			envValue:    "foobarbaz",
			wantTraceID: trace.TraceID{},
			wantSpanID:  trace.SpanID{},
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setting up the environment variable
			os.Setenv(tc.envKey, tc.envValue)

			// Call the function under test
			ctx, err := NewContextFromEnvTraceParent(context.Background())
			if tc.wantErr {
				assert.Error(t, err)
			} else if err != nil {
				t.Errorf("unexepcted error returned %s", err.Error())
			} else {
				spanCtx := trace.SpanContextFromContext(ctx)

				// If we do not expect an error, then the original & returned context should not be equal
				if !spanCtx.IsValid() || spanCtx.TraceID() != tc.wantTraceID || spanCtx.SpanID() != tc.wantSpanID {
					t.Errorf("unexpected span context")
				}
			}
		})

		// Cleanup after each test
		os.Unsetenv(tc.envKey)
	}
}

func TestNewContextFromDeterministicString(t *testing.T) {
	testCases := []struct {
		name          string
		traceIDString string
		spanIDString  string
		wantTraceID   trace.TraceID
		wantSpanID    trace.SpanID
		wantErr       bool
	}{
		{
			name:          "Test valid trace and span strings",
			traceIDString: "1234567",
			spanIDString:  "1234567",
			wantTraceID:   mustTraceIDFromHex("fcea920f7412b5da7be0cf42b8c93759"),
			wantSpanID:    mustSpanIDFromHex("fcea920f7412b5da"),
			wantErr:       false,
		},
		{
			name:          "Test empty trace and span strings",
			traceIDString: "",
			spanIDString:  "",
			wantTraceID:   mustTraceIDFromHex("d41d8cd98f00b204e9800998ecf8427e"),
			wantSpanID:    mustSpanIDFromHex("d41d8cd98f00b204"),
			wantErr:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := NewContextFromDeterministicString(tc.traceIDString, tc.spanIDString)
			span := trace.SpanContextFromContext(ctx)

			if tc.wantErr {
				if span.IsValid() {
					t.Errorf("expected invalid span context, got valid")
				}
			} else {
				if !span.IsValid() || span.TraceID() != tc.wantTraceID || span.SpanID() != tc.wantSpanID {
					t.Errorf("unexpected span context")
				}
			}
		})
	}
}

func TestGenTraceIDFromString(t *testing.T) {
	traceID, err := genTraceIDFromString("1234567")
	assert.Nil(t, err)
	assert.Equal(t, traceID, mustTraceIDFromHex("fcea920f7412b5da7be0cf42b8c93759"))
}

func TestGenSpanIDFromString(t *testing.T) {
	spanID, err := genSpanIDFromString("1234567")
	assert.Nil(t, err)
	assert.Equal(t, spanID, mustSpanIDFromHex("fcea920f7412b5da"))
}

func TestAttributeMapToKeyValue(t *testing.T) {
	testCases := []struct {
		name       string
		attributes map[string]string
		want       []attribute.KeyValue
	}{
		{
			name: "Single",
			attributes: map[string]string{
				"test": "value",
			},
			want: []attribute.KeyValue{
				attribute.String("test", "value"),
			},
		},
		{
			name: "Multiple",
			attributes: map[string]string{
				"test1": "value1",
				"test2": "value2",
			},
			want: []attribute.KeyValue{
				attribute.String("test1", "value1"),
				attribute.String("test2", "value2"),
			},
		},
		{
			name:       "Empty",
			attributes: map[string]string{},
			want:       nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := AttributeMapToKeyValue(tc.attributes); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AttributeMapToKeyValue() = %v, want %v", got, tc.want)
			}
		})
	}
}

func mustTraceIDFromHex(s string) (t trace.TraceID) {
	var err error
	t, err = trace.TraceIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}

func mustSpanIDFromHex(s string) (t trace.SpanID) {
	var err error
	t, err = trace.SpanIDFromHex(s)
	if err != nil {
		panic(err)
	}
	return
}
