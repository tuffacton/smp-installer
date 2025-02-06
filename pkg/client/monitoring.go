package client

import (
	"context"
	"path"

	"github.com/harness/smp-installer/pkg/profiles"
	"github.com/harness/smp-installer/pkg/render"
	"github.com/harness/smp-installer/pkg/store"
	"github.com/harness/smp-installer/pkg/tofu"
	"github.com/rs/zerolog/log"
)

type monitoringClient struct {
	clientConfig ClientConfig
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (m *monitoringClient) Exec(ctx context.Context) error {
	if !m.clientConfig.IsManaged {
		log.Info().Msgf("skipping %s sync as its not set to managed", m.clientConfig.ResourceName)
		return nil
	}
	tofu.ExecuteCommand(ctx, m.clientConfig.ContextDirectory, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, m.clientConfig.ContextDirectory, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (m *monitoringClient) PostExec(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// PreExec implements ResourceClient.
func (m *monitoringClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(m.clientConfig.Provider.Name(), m.clientConfig.ResourceName), m.clientConfig.ContextDirectory)
	if err != nil {
		return err
	}
	_, err = profiles.CopyFiles(m.clientConfig.ContextDirectory, m.clientConfig.Provider.Name(),
		[]string{"prometheus-override.yaml", "grafana-override.yaml"})
	if err != nil {
		log.Err(err).Msgf("unable to copy override files for monitoring module")
	}
	renderer := render.NewTemplateRenderer(m.configStore, m.outputStore)
	return renderer.Render(ctx, map[string]interface{}{}, m.clientConfig.ContextDirectory)
}

func NewMonitoringClient(clientConfig ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &monitoringClient{
		clientConfig: clientConfig,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}
