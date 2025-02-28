package resources

import (
	"context"

	"github.com/harness/smp-installer/pkg/client"
	"github.com/harness/smp-installer/pkg/store"
	"github.com/rs/zerolog/log"
)

type vpcCommand struct {
}

// Name implements ResourceCommand.
func (v *vpcCommand) Name() string {
	return "vpc"
}

// Sync implements ResourceCommand.
func (v *vpcCommand) Sync(ctx context.Context, configStore store.DataStore, outputStore store.DataStore) error {
	clientConfig := CreateClientConfig(ctx, v, configStore)
	vpcClient := client.NewVPCClient(clientConfig, configStore, outputStore)
	err := vpcClient.PreExec(ctx)
	if err != nil {
		log.Error().Msgf("pre-exec step failed while syncing %s", v.Name())
		return err
	}
	err = vpcClient.Exec(ctx)
	if err != nil {
		log.Error().Msgf("exec step failed while syncing %s", v.Name())
		return err
	}
	out, err := vpcClient.PostExec(ctx)
	if err != nil {
		log.Error().Msgf("post-exec step failed while syncing %s", v.Name())
		return err
	}
	outputStore.Set(ctx, v.Name(), out)
	return nil
}

func NewVpcCommand() ResourceCommand {
	return &vpcCommand{}
}
