package client

import (
	"context"
	"os"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/profiles"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type InfraCommand string

const (
	InitCommand  InfraCommand = "init"
	PlanCommand  InfraCommand = "plan"
	ApplyCommand InfraCommand = "apply"
)

type kubernetesClient struct {
	clientConfig ClientConfig
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (k *kubernetesClient) Exec(ctx context.Context) error {
	if !k.clientConfig.IsManaged {
		log.Info().Msgf("skipping %s sync as its not set to managed", k.clientConfig.ResourceName)
		return nil
	}
	tofu.ExecuteCommand(ctx, k.clientConfig.ContextDirectory, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, k.clientConfig.ContextDirectory, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (k *kubernetesClient) PostExec(ctx context.Context) (map[string]interface{}, error) {
	clusterName := ""
	if !k.clientConfig.IsManaged {
		clusterName = k.configStore.GetString(ctx, "cluster_name")
	} else {
		var err error = nil
		clusterName, err = tofu.GetOutput(ctx, k.clientConfig.ContextDirectory, "clustername", ".")
		if err != nil {
			log.Err(err).Msgf("unable to retrieve cluster name")
			return nil, err
		}
	}
	return map[string]interface{}{"cluster_name": clusterName}, nil
}

// PreExec implements ResourceClient.
func (k *kubernetesClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(k.clientConfig.Provider.Name(), k.clientConfig.ResourceName), k.clientConfig.ContextDirectory)
	if err != nil {
		return err
	}
	profile := k.configStore.GetString(ctx, store.ProfileKey)
	err = profiles.CopyFiles(profile, k.clientConfig.ContextDirectory)
	if err != nil {
		log.Err(err).Msgf("error copying profile files")
		return err
	}
	profileConfig := make(map[string]interface{})
	configData, err := os.ReadFile(path.Join(k.clientConfig.ContextDirectory, "config.yaml"))
	if err != nil {
		log.Err(err).Msgf("unable to read profile config file")
		return err
	}
	err = yaml.Unmarshal(configData, &profileConfig)
	if err != nil {
		log.Err(err).Msgf("unable to unmarshal yaml data from profile config")
		return err
	}
	err = k.configStore.AddAll(ctx, profileConfig)
	renderer := render.NewTemplateRenderer(k.configStore, k.outputStore)
	return renderer.Render(ctx, map[string]interface{}{}, k.clientConfig.ContextDirectory)
}

func NewKubernetesClient(clientConfig ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &kubernetesClient{
		clientConfig: clientConfig,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}
