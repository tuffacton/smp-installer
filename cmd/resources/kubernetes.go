package resources

import (
	"context"

	"github.com/harness/smp-installer/pkg/client"
	"github.com/harness/smp-installer/pkg/store"
	"github.com/rs/zerolog/log"
)

type kubernetesCommand struct {
}

// Name implements ResourceCommand.
func (l *kubernetesCommand) Name() string {
	return "kubernetes"
}

// Sync implements ResourceCommand.
func (l *kubernetesCommand) Sync(ctx context.Context, configStore store.DataStore, outputStore store.DataStore) error {
	clientConfig := CreateClientConfig(ctx, l, configStore)
	lbClient := client.NewKubernetesClient(clientConfig, configStore, outputStore)
	err := lbClient.PreExec(ctx)
	if err != nil {
		log.Error().Msgf("pre-exec step failed while syncing %s", l.Name())
		return err
	}
	err = lbClient.Exec(ctx)
	if err != nil {
		log.Error().Msgf("exec step failed while syncing %s", l.Name())
		return err
	}
	out, err := lbClient.PostExec(ctx)
	if err != nil {
		log.Error().Msgf("post-exec step failed while syncing %s", l.Name())
		return err
	}
	outputStore.Set(ctx, l.Name(), out)
	return nil
}

func NewKubernetesCommand() ResourceCommand {
	return &kubernetesCommand{}
}
