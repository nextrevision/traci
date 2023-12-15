package cmd

import (
	"context"
	"fmt"
	"github.com/nextrevision/traci/providers"
	"github.com/nextrevision/traci/tracing"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "execute the command provided",
	Long: `execute the command provided as an otel span, measuring and send the result to the otel server.
The wrapping span's w3c traceparent is automatically passed to the child process's environment as TRACEPARENT.

Examples:

traci exec curl https://duckduckgo.com

traci exec /bin/sh -c 'traci exec curl https://duckduckgo.com && sleep 1'

TRACI_SERVICE_NAME=foo traci exec curl https://duckduckgo.com`,
	RunE: runCommand,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	// Disable printing usage due to the return type always being a non-nil error
	execCmd.SilenceUsage = true
	// Handle errors ourselves
	execCmd.SilenceErrors = true
	// Do not parse flags for the command and pass all flags to the child exec process
	execCmd.DisableFlagParsing = true

	rootCmd.AddCommand(execCmd)
}

func runCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	config := getConfig()

	ciProvider := providers.DetectProvider()

	serviceName := ciProvider.GetServiceName()
	if config.ServiceName != "" {
		serviceName = config.ServiceName
	}

	spanName := ciProvider.GetSpanName()
	if config.SpanName != "" {
		spanName = config.SpanName
	}

	command := args[0]
	commandPath, _ := exec.LookPath(command)

	// Set resource attributes on the span
	var resourceAttributes []attribute.KeyValue
	resourceAttributes = append(resourceAttributes, semconv.ProcessExecutableName(command))
	resourceAttributes = append(resourceAttributes, semconv.ProcessExecutablePath(commandPath))
	resourceAttributes = append(resourceAttributes, tracing.AttributeMapToKeyValue(ciProvider.GetAttributes())...)
	resourceAttributes = append(resourceAttributes, attribute.String("traci.ci.provider", ciProvider.GetCIName()))
	resourceAttributes = append(resourceAttributes, attribute.String("traci.version", rootCmd.Version))

	// Add command args as a process attribute to the span if specified
	if config.TagCommandArgs && len(args) > 1 {
		resourceAttributes = append(resourceAttributes, semconv.ProcessCommandArgs(args[1:]...))
	}

	traceCtx, err := tracing.NewContextFromEnvTraceParent(ctx)
	if err != nil {
		traceCtx = tracing.NewContextFromDeterministicString(ciProvider.GetTraceVal(), ciProvider.GetSpanVal())
	}

	traceProvider := tracing.NewTraceProvider(traceCtx, serviceName, resourceAttributes)

	tracer := tracing.NewTracer(serviceName, traceProvider)

	// Start a trace
	spanCtx, span := tracer.Start(traceCtx, spanName)

	var child *exec.Cmd
	if len(args) > 1 {
		child = exec.CommandContext(spanCtx, args[0], args[1:]...)
	} else {
		child = exec.CommandContext(spanCtx, args[0])
	}

	// attach all stdio to the parent's handles
	child.Stdin = cmd.InOrStdin()
	child.Stdout = cmd.OutOrStdout()
	child.Stderr = cmd.ErrOrStderr()

	// Replace TRACEPARENT in the environment with one from this span
	child.Env = []string{}
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, fmt.Sprintf("%s=", tracing.TraceParentKey)) {
			child.Env = append(child.Env, env)
		}
	}
	child.Env = append(child.Env, fmt.Sprintf("%s=%s", tracing.TraceParentKey, tracing.GenTraceParentString(span.SpanContext())))

	// Forward CTRL-C (SIGINT) to the child process to attempt graceful shutdown
	signals := make(chan os.Signal, 10)
	signalsDone := make(chan struct{})
	signal.Notify(signals, os.Interrupt)
	go func() {
		sig := <-signals
		child.Process.Signal(sig)
		close(signalsDone)
	}()

	// Run the child process and record any errors
	if err = child.Run(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		slog.Debug(err.Error())
	}

	errCode := ErrorCode{
		Code: child.ProcessState.ExitCode(),
		Err:  err,
	}

	// Send the span to the collector and force a shutdown of the TraceProvider with a timeout
	sent := make(chan bool, 1)
	go func() {
		span.End()

		ctxTimeout, cancel := context.WithTimeout(ctx, time.Millisecond*100)
		defer cancel()

		err := traceProvider.ForceFlush(ctxTimeout)
		if err != nil {
			slog.Debug(err.Error())
		}
		err = traceProvider.Shutdown(ctxTimeout)
		if err != nil {
			slog.Debug(err.Error())
		}
		sent <- true
	}()

	select {
	case <-sent:
		return &errCode
	case <-time.After(500 * time.Millisecond): // TODO Make configurable
		return &errCode
	}
}
