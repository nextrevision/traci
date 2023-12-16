package config

type Config struct {
	ServiceName    string `mapstructure:"service_name"`
	SpanName       string `mapstructure:"span_name"`
	TraceBoundary  string `mapstructure:"trace_boundary" default:"pipeline"`
	TagCommandArgs bool   `mapstructure:"tag_command_args"`
}

type TraceBoundary string

const (
	TraceBoundaryPipeline TraceBoundary = "pipeline"
	TraceBoundaryJob      TraceBoundary = "job"
)
