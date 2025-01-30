package resources

import (
	"context"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/client"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"github.com/rs/zerolog/log"
)

type secretOperatorCommand struct {
}

// Name implements ResourceCommand.
func (l *secretOperatorCommand) Name() string {
	return "secretoperator"
}

// Sync implements ResourceCommand.
func (l *secretOperatorCommand) Sync(ctx context.Context, configStore store.DataStore, outputStore store.DataStore) error {
	clientConfig := CreateClientConfig(ctx, l, configStore)
	operatorClient := client.NewSecretOperatorClient(clientConfig, configStore, outputStore)
	err := operatorClient.PreExec(ctx)
	if err != nil {
		log.Error().Msgf("pre-exec step failed while syncing %s", l.Name())
		return err
	}
	err = operatorClient.Exec(ctx)
	if err != nil {
		log.Error().Msgf("exec step failed while syncing %s", l.Name())
		return err
	}
	out, err := operatorClient.PostExec(ctx)
	if err != nil {
		log.Error().Msgf("post-exec step failed while syncing %s", l.Name())
		return err
	}
	outputStore.Set(ctx, l.Name(), out)
	return nil
}

func NewSecretOperatorCommand() ResourceCommand {
	return &secretOperatorCommand{}
}
