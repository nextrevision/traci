package providers

import (
	"os"
	"strings"
)

type Travis struct{}

func (t Travis) GetCIName() string {
	return "Travis-CI"
}

func (t Travis) GetTraceVal() string {
	return os.Getenv("TRAVIS_BUILD_ID")
}

func (t Travis) GetSpanVal() string {
	return os.Getenv("TRAVIS_JOB_ID")
}

func (t Travis) GetServiceName() string {
	return strings.ToLower(t.GetCIName())
}

func (t Travis) GetSpanName() string {
	return os.Getenv("TRAVIS_JOB_NAME")
}

func (t Travis) GetAttributes() map[string]string {
	// See https://docs.travis-ci.com/user/environment-variables/
	return map[string]string{
		"travis.build.branch": os.Getenv("TRAVIS_BRANCH"),
		"travis.build.id":     os.Getenv("TRAVIS_BUILD_ID"),
		"travis.build.url":    os.Getenv("TRAVIS_BUILD_WEB_URL"),
		"travis.job.name":     os.Getenv("TRAVIS_JOB_NAME"),
		"travis.job.number":   os.Getenv("TRAVIS_JOB_NUMBER"),
		"travis.job.url":      os.Getenv("TRAVIS_JOB_WEB_URL"),
		"travis.repo":         os.Getenv("TRAVIS_REPO_SLUG"),
		"travis.sha":          os.Getenv("TRAVIS_COMMIT"),
	}

}
