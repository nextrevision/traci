package providers

import (
	"os"
	"strings"
)

type Bitbucket struct{}

func (b Bitbucket) GetCIName() string {
	return "Bitbucket"
}

func (b Bitbucket) GetTraceString() string {
	return os.Getenv("BITBUCKET_PIPELINE_UUID")
}

func (b Bitbucket) GetServiceName() string {
	return strings.ToLower(b.GetCIName())
}

func (b Bitbucket) GetSpanName() string {
	// Bitbucket doesn't support a good variable for step name; if they ever do, this should change
	return os.Getenv("BITBUCKET_STEP_UUID")
}

func (b Bitbucket) GetAttributes() map[string]string {
	// See https://support.atlassian.com/bitbucket-cloud/docs/variables-and-secrets/
	return map[string]string{
		"bitbucket.build":     os.Getenv("BITBUCKET_BUILD_NUMBER"),
		"bitbucket.sha":       os.Getenv("BITBUCKET_COMMIT"),
		"bitbucket.branch":    os.Getenv("BITBUCKET_BRANCH"),
		"bitbucket.repo.slug": os.Getenv("BITBUCKET_REPO_SLUG"),
		"bitbucket.repo.uuid": os.Getenv("BITBUCKET_REPO_UUID"),
	}
}
