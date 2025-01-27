package client

import (
	"context"
)

type CloudProviderType string

const (
	AWSCloudProvider CloudProviderType = "aws"
	GCPCloudProvider CloudProviderType = "gcp"
)

func (cpt CloudProviderType) Name() string {
	return string(cpt)
}

type ClientConfig struct {
	OutputDirectory  string
	Provider         CloudProviderType
	IsManaged        bool
	ResourceName     string
	ContextDirectory string
}

type ResourceClient interface {
	PreExec(ctx context.Context) error
	Exec(ctx context.Context) error
	PostExec(ctx context.Context) (map[string]interface{}, error)
}
