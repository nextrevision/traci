package providers

import (
	"os"
	"strings"
)

type CircleCI struct{}

func (c CircleCI) GetCIName() string {
	return "CircleCI"
}

func (c CircleCI) GetTraceVal() string {
	return os.Getenv("CIRCLE_WORKFLOW_ID")
}

func (c CircleCI) GetSpanVal() string {
	return os.Getenv("CIRCLE_WORKFLOW_JOB_ID")
}

func (c CircleCI) GetServiceName() string {
	return strings.ToLower(c.GetCIName())
}

func (c CircleCI) GetSpanName() string {
	return os.Getenv("CIRCLE_JOB")
}

func (c CircleCI) GetAttributes() map[string]string {
	// See https://circleci.com/docs/variables/
	return map[string]string{
		"circleci.project.name":    os.Getenv("CIRCLE_PROJECT_REPONAME"),
		"circleci.workflow.id":     os.Getenv("CIRCLE_WORKFLOW_ID"),
		"circleci.workflow.job_id": os.Getenv("CIRCLE_WORKFLOW_JOB_ID"),
		"circleci.build.num":       os.Getenv("CIRCLE_BUILD_NUM"),
		"circleci.build.url":       os.Getenv("CIRCLE_BUILD_URL"),
		"circleci.job":             os.Getenv("CIRCLE_JOB"),
		"circleci.sha":             os.Getenv("CIRCLE_SHA1"),
	}
}
