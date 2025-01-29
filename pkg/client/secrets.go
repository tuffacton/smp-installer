package client

import (
	"context"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
)

type secretsClient struct {
	clientConfig ClientConfig
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (s *secretsClient) Exec(ctx context.Context) error {
	if !s.clientConfig.IsManaged {
		return nil
	}
	tofu.ExecuteCommand(ctx, s.clientConfig.ContextDirectory, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, s.clientConfig.ContextDirectory, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (s *secretsClient) PostExec(ctx context.Context) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// PreExec implements ResourceClient.
func (s *secretsClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(s.clientConfig.Provider.Name(),
		s.clientConfig.ResourceName),
		s.clientConfig.ContextDirectory)
	if err != nil {
		return err
	}
	renderer := render.NewTemplateRenderer(s.configStore, s.outputStore)
	return renderer.Render(ctx, make(map[string]interface{}), s.clientConfig.ContextDirectory)
}

func NewSecretsClient(clientConfig ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &secretsClient{
		clientConfig: clientConfig,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}
