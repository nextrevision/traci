package cmd

import (
	"github.com/nextrevision/traci/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var execfCmd = &cobra.Command{
	Use:   "execf",
	Short: "execute the command provided with support for traci cli flags",
	Long: `execute the command provided as an otel span, measuring and send the result to the otel server.
	The wrapping span's w3c traceparent is automatically passed to the child process's environment as TRACEPARENT.

Examples:

traci execf --service-name foo -- curl https://duckduckgo.com

traci execf --span-name bar -- /bin/sh -c 'traci exec curl https://duckduckgo.com && sleep 1'`,
	RunE: runCommand,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	// Disable printing usage due to the return type always being a non-nil error
	execfCmd.SilenceUsage = true
	// Handle errors ourselves
	execfCmd.SilenceErrors = true

	var traceBoundaryValue = &EnumValue{
		Allowed: []string{
			string(config.TraceBoundaryPipeline),
			string(config.TraceBoundaryJob),
		},
	}

	execfCmd.Flags().StringP("span-name", "s", "", "name of the span")
	execfCmd.Flags().StringP("service-name", "n", "", "name of the service")
	execfCmd.Flags().VarP(traceBoundaryValue, "trace-boundary", "t", "limit the trace to a pipeline, stage or job")
	execfCmd.Flags().Bool("tag-command-args", false, "tag spans with the full list of command arguments")

	viper.BindPFlag("span_name", execfCmd.Flags().Lookup("span-name"))
	viper.BindPFlag("service_name", execfCmd.Flags().Lookup("service-name"))
	viper.BindPFlag("trace_boundary", execfCmd.Flags().Lookup("trace-boundary"))
	viper.BindPFlag("tag_command_args", execfCmd.Flags().Lookup("tag-command-args"))

	rootCmd.AddCommand(execfCmd)
}
