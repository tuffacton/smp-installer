package resources

import (
	"context"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/client"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
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
	clientConfig, err := CreateClientConfig(ctx, l, configStore)
	if err != nil {
		log.Err(err).Msgf("unable to create client config")
		return err
	}
	helmClient := client.NewHelmClient(clientConfig, configStore, outputStore)
	err = helmClient.PreExec(ctx)
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
