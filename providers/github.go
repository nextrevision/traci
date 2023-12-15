package providers

import (
	"fmt"
	"os"
	"strings"
)

type GitHubActions struct{}

func (g GitHubActions) GetCIName() string {
	return "GitHub-Actions"
}

func (g GitHubActions) GetTraceString() string {
	return fmt.Sprintf("%s-%s-%s", os.Getenv("GITHUB_RUN_ID"), os.Getenv("GITHUB_RUN_NUMBER"), os.Getenv("GITHUB_RUN_ATTEMPT"))
}

func (g GitHubActions) GetServiceName() string {
	return strings.ToLower(g.GetCIName())
}

func (g GitHubActions) GetSpanName() string {
	return os.Getenv("GITHUB_JOB")
}

func (g GitHubActions) GetAttributes() map[string]string {
	// See https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
	return map[string]string{
		"github.action":        os.Getenv("GITHUB_ACTION"),
		"github.action.repo":   os.Getenv("GITHUB_ACTION_REPOSITORY"),
		"github.workflow":      os.Getenv("GITHUB_WORKFLOW"),
		"github.job.id":        os.Getenv("GITHUB_JOB"),
		"github.run.id":        os.Getenv("GITHUB_RUN_ID"),
		"github.run.number":    os.Getenv("GITHUB_RUN_NUMBER"),
		"github.ref.name":      os.Getenv("GITHUB_REF"),
		"github.ref.type":      os.Getenv("GITHUB_REF_TYPE"),
		"github.repo.name":     os.Getenv("GITHUB_REPOSITORY"),
		"github.repo.id":       os.Getenv("GITHUB_REPOSITORY_ID"),
		"github.repo.owner":    os.Getenv("GITHUB_REPOSITORY_OWNER"),
		"github.repo.owner-id": os.Getenv("GITHUB_REPOSITORY_OWNER_ID"),
		"github.sha":           os.Getenv("GITHUB_SHA"),
	}
}
