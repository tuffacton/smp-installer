package client

import (
	"context"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"github.com/rs/zerolog/log"
)

type secretOperatorClient struct {
	clientConfig ClientConfig
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (s *secretOperatorClient) Exec(ctx context.Context) error {
	if !s.clientConfig.IsManaged {
		return nil
	}
	tofu.ExecuteCommand(ctx, s.clientConfig.ContextDirectory, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, s.clientConfig.ContextDirectory, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (s *secretOperatorClient) PostExec(ctx context.Context) (map[string]interface{}, error) {
	if !s.clientConfig.IsManaged {
		serviceAccount := s.configStore.GetString(ctx, "secretoperator.kubernetes_service_account")
		return map[string]interface{}{"service_account": serviceAccount}, nil
	}
	serviceAccount, err := tofu.GetOutput(ctx, s.clientConfig.ContextDirectory, "service_account", ".")
	if err != nil {
		log.Err(err).Msgf("unable to retrieve service account")
		return nil, err
	}
	return map[string]interface{}{"service_account": serviceAccount}, nil
}

// PreExec implements ResourceClient.
func (s *secretOperatorClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(s.clientConfig.Provider.Name(),
		s.clientConfig.ResourceName),
		s.clientConfig.ContextDirectory)
	if err != nil {
		return err
	}
	renderer := render.NewTemplateRenderer(s.configStore, s.outputStore)
	return renderer.Render(ctx, make(map[string]interface{}), s.clientConfig.ContextDirectory)
}

func NewSecretOperatorClient(clientConfig ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &secretOperatorClient{
		clientConfig: clientConfig,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}
