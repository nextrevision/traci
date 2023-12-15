package config

type Config struct {
	ServiceName    string `mapstructure:"service_name"`
	SpanName       string `mapstructure:"span_name"`
	TagCommandArgs bool   `mapstructure:"tag_command_args"`
}
