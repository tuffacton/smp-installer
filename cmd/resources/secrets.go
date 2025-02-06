package resources

import (
	"context"

	"github.com/harness/smp-installer/pkg/client"
	"github.com/harness/smp-installer/pkg/store"
	"github.com/rs/zerolog/log"
)

type secretsCommand struct {
}

// Name implements ResourceCommand.
func (l *secretsCommand) Name() string {
	return "secrets"
}

// Sync implements ResourceCommand.
func (l *secretsCommand) Sync(ctx context.Context, configStore store.DataStore, outputStore store.DataStore) error {
	clientConfig := CreateClientConfig(ctx, l, configStore)
	helmClient := client.NewSecretsClient(clientConfig, configStore, outputStore)
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

func NewSecretsCommand() ResourceCommand {
	return &secretsCommand{}
}
