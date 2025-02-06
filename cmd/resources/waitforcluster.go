package resources

import (
	"context"

	"github.com/harness/smp-installer/pkg/client"
	"github.com/harness/smp-installer/pkg/store"
	"github.com/rs/zerolog/log"
)

type waitForClusterCommand struct {
}

// Name implements ResourceCommand.
func (w *waitForClusterCommand) Name() string {
	return "waitforcluster"
}

// Sync implements ResourceCommand.
func (w *waitForClusterCommand) Sync(ctx context.Context,
	configStore store.DataStore, outputStore store.DataStore) error {
	clientConfig := CreateClientConfig(ctx, w, configStore)
	waitForClusterClient := client.NewWaitForClusterClient(clientConfig, configStore, outputStore)
	err := waitForClusterClient.PreExec(ctx)
	if err != nil {
		log.Error().Msgf("pre-exec step failed while syncing %s", w.Name())
		return err
	}
	err = waitForClusterClient.Exec(ctx)
	if err != nil {
		log.Error().Msgf("exec step failed while syncing %s", w.Name())
		return err
	}
	out, err := waitForClusterClient.PostExec(ctx)
	if err != nil {
		log.Error().Msgf("post-exec step failed while syncing %s", w.Name())
		return err
	}
	outputStore.Set(ctx, w.Name(), out)
	return nil
}

func NewWaitForClusterCommand() ResourceCommand {
	return &waitForClusterCommand{}
}
