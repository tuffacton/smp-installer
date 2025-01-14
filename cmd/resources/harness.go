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
	helmClient := client.NewHelmClient(l.Name(), configStore, outputStore)
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
	err = helmClient.PostExec(ctx)
	if err != nil {
		log.Error().Msgf("post-exec step failed while syncing %s", l.Name())
		return err
	}
	return nil
}

func NewHarnessCommand() ResourceCommand {
	return &harnessCommand{}
}
