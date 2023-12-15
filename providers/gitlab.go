package providers

import (
	"os"
	"strings"
)

type GitLabCI struct{}

func (g GitLabCI) GetCIName() string {
	return "GitLab-CI"
}

func (g GitLabCI) GetTraceString() string {
	return os.Getenv("CI_PIPELINE_ID")
}

func (g GitLabCI) GetServiceName() string {
	return strings.ToLower(g.GetCIName())
}

func (g GitLabCI) GetSpanName() string {
	return os.Getenv("CI_JOB_NAME")
}

func (g GitLabCI) GetAttributes() map[string]string {
	// See https://docs.gitlab.com/ee/ci/variables/predefined_variables.html
	return map[string]string{
		"gitlab.project.id":   os.Getenv("CI_PROJECT_ID"),
		"gitlab.project.name": os.Getenv("CI_PROJECT_NAME"),
		"gitlab.pipeline.id":  os.Getenv("CI_PIPELINE_ID"),
		"gitlab.pipeline.ref": os.Getenv("CI_COMMIT_REF_NAME"),
		"gitlab.pipeline.sha": os.Getenv("CI_COMMIT_SHA"),
		"gitlab.job.id":       os.Getenv("CI_JOB_ID"),
		"gitlab.job.name":     os.Getenv("CI_JOB_NAME"),
		"gitlab.job.stage":    os.Getenv("CI_JOB_STAGE"),
		"gitlab.job.ref":      os.Getenv("CI_COMMIT_REF_NAME"),
		"gitlab.job.sha":      os.Getenv("CI_COMMIT_SHA"),
	}
}
