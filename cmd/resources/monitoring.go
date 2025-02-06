package resources

import (
	"context"

	"github.com/harness/smp-installer/pkg/client"
	"github.com/harness/smp-installer/pkg/store"
	"github.com/rs/zerolog/log"
)

type monitoringCommand struct {
}

// Name implements ResourceCommand.
func (m *monitoringCommand) Name() string {
	return "monitoring"
}

// Sync implements ResourceCommand.
func (m *monitoringCommand) Sync(ctx context.Context, configStore store.DataStore, outputStore store.DataStore) error {
	clientConfig := CreateClientConfig(ctx, m, configStore)
	lbClient := client.NewLoadbalancerClient(clientConfig, configStore, outputStore)
	err := lbClient.PreExec(ctx)
	if err != nil {
		log.Error().Msgf("pre-exec step failed while syncing %s", m.Name())
		return err
	}
	err = lbClient.Exec(ctx)
	if err != nil {
		log.Error().Msgf("exec step failed while syncing %s", m.Name())
		return err
	}
	out, err := lbClient.PostExec(ctx)
	if err != nil {
		log.Error().Msgf("post-exec step failed while syncing %s", m.Name())
		return err
	}
	outputStore.Set(ctx, m.Name(), out)
	return nil
}

func NewMonitoringCommand() ResourceCommand {
	return &monitoringCommand{}
}
