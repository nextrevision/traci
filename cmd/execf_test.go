package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecfCmd(t *testing.T) {
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
			args:    []string{"execf", "echo"},
			stdout:  "",
			stderr:  "",
			rc:      0,
			wantErr: false,
		},
		{
			name:    "command with args after end of cmd options",
			args:    []string{"execf", "--", "echo", "-e", "foobar"},
			stdout:  "foobar",
			stderr:  "",
			rc:      0,
			wantErr: false,
		},
		{
			name:    "command with flag",
			args:    []string{"execf", "--span-name", "foobar", "--", "echo", "foobar"},
			stdout:  "foobar",
			stderr:  "",
			rc:      0,
			wantErr: false,
		},
		{
			name:    "command with invalid flag",
			args:    []string{"execf", "--test", "--", "echo", "foobar"},
			stdout:  "",
			stderr:  "",
			rc:      1,
			wantErr: true,
		},
	}

	rootCmd.AddCommand(execfCmd)

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
