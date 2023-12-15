package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "traci",
	Short:   "traci wraps commands in OpenTelemetry tracing and sends them to a OTel endpoint",
	Run:     func(cmd *cobra.Command, args []string) {},
	Version: "0.1.0",
}

func Execute() {
	var e *ErrorCode
	res := rootCmd.Execute()

	// Handle error codes if passed from the child command
	if errors.As(res, &e) {
		if e.Code != 0 && e.Err != nil {
			switch e.Err.(type) {
			case *exec.ExitError:
				// Do not print redundant exit status error messages
			case *fs.PathError:
				// Removing the "fork/exec" prefix from command not found messages
				fmt.Fprintln(os.Stderr, strings.Replace(e.Err.Error(), "fork/exec ", "", 1))
			default:
				// Print all other errors
				fmt.Fprintln(os.Stderr, e.Err.Error())
			}
		}
		os.Exit(e.Code)
	} else if res != nil {
		fmt.Fprintln(os.Stderr, res.Error())
		os.Exit(1)
	}
}

type ErrorCode struct {
	Code int
	Err  error
}

func NewErrorCode(code int, err error) *ErrorCode {
	if code > 0 && err == nil {
		err = errors.New("non-zero exit code from command")
	}
	return &ErrorCode{Code: code, Err: err}
}

func (e *ErrorCode) Error() string {
	if e.Err == nil {
		return ""
	}
	return e.Err.Error()
}

func init() {
	viper.SetEnvPrefix("traci")
	viper.AutomaticEnv()
}
