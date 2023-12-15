package providers

import (
	"os"
	"reflect"
	"testing"
)

// Create a test that will fail if the provider is not detected.
func TestDetectProvider(t *testing.T) {
	tests := []struct {
		name   string
		envVar string
		want   Provider
	}{
		{"gitlab", "GITLAB_CI", &GitLabCI{}},
		{"circleci", "CIRCLECI", &CircleCI{}},
		{"travis", "TRAVIS", &Travis{}},
		{"github", "GITHUB_ACTION", &GitHubActions{}},
		{"bitbucket", "BITBUCKET_BUILD_NUMBER", &Bitbucket{}},
		{"default", "FOOBAR", &DefaultProvider{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer os.Clearenv()
			os.Setenv(tc.envVar, "true")

			got := DetectProvider()

			if reflect.TypeOf(got).Kind() == reflect.TypeOf(tc.want).Kind() {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
