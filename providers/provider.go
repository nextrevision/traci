package providers

import (
	"crypto/rand"
	"encoding/hex"
	"os"
)

func DetectProvider() Provider {
	if _, present := os.LookupEnv("GITLAB_CI"); present {
		return GitLabCI{}
	} else if _, present := os.LookupEnv("CIRCLECI"); present {
		return CircleCI{}
	} else if _, present := os.LookupEnv("TRAVIS"); present {
		return Travis{}
	} else if _, present := os.LookupEnv("GITHUB_ACTION"); present {
		return GitHubActions{}
	} else if _, present := os.LookupEnv("BITBUCKET_BUILD_NUMBER"); present {
		return Bitbucket{}
	}

	return DefaultProvider{}
}

type Provider interface {
	GetCIName() string
	GetPipelineID() string
	GetJobID() string
	GetServiceName() string
	GetSpanName() string
	GetAttributes() map[string]string
}

type DefaultProvider struct{}

func (d DefaultProvider) GetCIName() string {
	return "Default"
}

func (d DefaultProvider) GetPipelineID() string {
	return d.genTraceID()
}

func (d DefaultProvider) GetJobID() string {
	return d.genTraceID()
}

func (d DefaultProvider) GetServiceName() string {
	return "traci"
}

func (d DefaultProvider) GetSpanName() string {
	return "cmd"
}

func (d DefaultProvider) GetAttributes() map[string]string {
	return map[string]string{}
}

func (d DefaultProvider) genTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
