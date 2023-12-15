package cmd

import (
	"bytes"
	"errors"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestExecCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		stdout  string
		stderr  string
		rc      int
		wantErr bool
	}{
		{
			name:    "single command",
			args:    []string{"exec", "echo"},
			stdout:  "",
			stderr:  "",
			rc:      0,
			wantErr: false,
		},
		{
			name:    "command with args",
			args:    []string{"exec", "echo", "-e", "foobar"},
			stdout:  "foobar",
			stderr:  "",
			rc:      0,
			wantErr: false,
		},
		{
			name:    "invalid flag as command",
			args:    []string{"exec", "--span-name"},
			stdout:  "",
			stderr:  "",
			rc:      -1,
			wantErr: true,
		},
		{
			name:    "command with stderr out",
			args:    []string{"exec", "/bin/sh", "-c", "echo foobar 1>&2"},
			stdout:  "",
			stderr:  "foobar",
			rc:      0,
			wantErr: false,
		},
		{
			name:    "command with return code",
			args:    []string{"exec", "/bin/sh", "-c", "exit 3"},
			stdout:  "",
			stderr:  "",
			rc:      3,
			wantErr: true,
		},
	}

	rootCmd.AddCommand(execCmd)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, errCode := execute(t, rootCmd, tc.args...)
			if tc.wantErr {
				assert.NotNil(t, errCode.Err)
			} else {
				assert.Nil(t, errCode.Err)
			}

			assert.Equal(t, tc.rc, errCode.Code)

			assert.Equal(t, tc.stdout, stdout)
			assert.Equal(t, tc.stderr, stderr)
		})
	}
}

func execute(t *testing.T, c *cobra.Command, args ...string) (string, string, *ErrorCode) {
	t.Helper()

	bufOut := new(bytes.Buffer)
	bufErr := new(bytes.Buffer)
	c.SetOut(bufOut)
	c.SetErr(bufErr)
	c.SetArgs(args)

	err := c.Execute()

	var e *ErrorCode
	stdout := strings.TrimSpace(bufOut.String())
	stderr := strings.TrimSpace(bufErr.String())

	if errors.As(err, &e) {
		return stdout, stderr, e
	}

	if err != nil {
		return stdout, stderr, NewErrorCode(1, err)
	}

	return stdout, stderr, NewErrorCode(0, err)
}
