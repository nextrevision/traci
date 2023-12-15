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
	GetTraceVal() string
	GetSpanVal() string
	GetServiceName() string
	GetSpanName() string
	GetAttributes() map[string]string
}

type DefaultProvider struct{}

func (d DefaultProvider) GetCIName() string {
	return "Default"
}

func (d DefaultProvider) GetTraceVal() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func (d DefaultProvider) GetSpanVal() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func (d DefaultProvider) GetServiceName() string {
	return "traci"
}

func (d DefaultProvider) GetSpanName() string {
	return ""
}

func (d DefaultProvider) GetAttributes() map[string]string {
	return map[string]string{}
}
