package client

import (
	"context"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"github.com/rs/zerolog/log"
)

type Configuration struct {
	OutputDirectory string               `yaml:"output_dir"`
	InfraConfig     InfraConfiguration   `yaml:"infra"`
	HarnessConfig   HarnessConfiguration `yaml:"harness"`
}

type CloudProviderType string

const (
	AWSCloudProvider CloudProviderType = "aws"
	GCPCloudProvider CloudProviderType = "gcp"
)

type ResourceClient interface {
	PreExec(ctx context.Context) error
	Exec(ctx context.Context) error
	PostExec(ctx context.Context) error
}

func CheckIsManaged(ctx context.Context, resourceName string, configStore store.DataStore) (bool, error) {
	isManaged, err := configStore.Get(ctx, resourceName+".manage")
	if err != nil {
		log.Err(err).Msgf("config store is wrongly configured")
		return false, err
	}
	isManagedB := isManaged.(bool)
	return isManagedB, nil
}
