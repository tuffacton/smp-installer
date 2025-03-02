package resources

import (
	"context"

	"github.com/harness/smp-installer/pkg/client"
	"github.com/harness/smp-installer/pkg/store"
	"github.com/rs/zerolog/log"
)

type harnessCommand struct {
}

// Name implements ResourceCommand.
func (l *harnessCommand) Name() string {
	return "harness"
}

// Sync implements ResourceCommand.
func (l *harnessCommand) Sync(ctx context.Context, configStore store.DataStore, outputStore store.DataStore) error {
	// externalSecretsEnabled := configStore.GetBool(context.Background(), "secrets.manage")
	clientConfig := CreateClientConfig(ctx, l, configStore)
	helmClient := client.NewHelmClient(clientConfig, configStore, outputStore)
	err := helmClient.PreExec(ctx)
	if err != nil {
		log.Error().Msgf("pre-exec step failed while syncing %s", l.Name())
		return err
	}
	err = helmClient.Exec(ctx)
	if err != nil {
		log.Error().Msgf("exec step failed while syncing %s", l.Name())
		return err
	}
	out, err := helmClient.PostExec(ctx)
	if err != nil {
		log.Error().Msgf("post-exec step failed while syncing %s", l.Name())
		return err
	}
	outputStore.Set(ctx, l.Name(), out)
	return nil
}

func NewHarnessCommand() ResourceCommand {
	return &harnessCommand{}
}
