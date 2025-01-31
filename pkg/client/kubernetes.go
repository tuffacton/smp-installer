package client

import (
	"context"
	"encoding/json"
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
func (k *kubernetesClient) PostExec(ctx context.Context) (output map[string]interface{}, err error) {
	clusterName := ""
	vpc := ""
	subnets := []string{}
	oidcProvider := ""
	oidcIssuer := ""
	if !k.clientConfig.IsManaged {
		clusterName = k.configStore.GetString(ctx, "kubernetes.cluster_name")
		vpc = k.configStore.GetString(ctx, "kubernetes.vpc")
		subnetsFromConfig, err := k.configStore.Get(ctx, "kubernetes.subnets")
		if err == nil {
			subnets = subnetsFromConfig.([]string)
		}
		oidcProvider = k.configStore.GetString(ctx, "kubernetes.oidc_provider_arn")
		oidcIssuer = k.configStore.GetString(ctx, "kubernetes.oidc_provider_url")
	} else {
		var err error = nil
		clusterName, err = tofu.GetOutput(ctx, k.clientConfig.ContextDirectory, "clustername", ".")
		if err != nil {
			log.Err(err).Msgf("unable to retrieve cluster name")
			return nil, err
		}
		vpc, err = tofu.GetOutput(ctx, k.clientConfig.ContextDirectory, "vpc", ".")
		if err != nil {
			log.Err(err).Msgf("unable to retrieve vpc name from tofu output")
			return nil, err
		}
		subnetsFromTofuOutput, err := tofu.GetOutput(ctx, k.clientConfig.ContextDirectory, "subnets", ".")
		if err != nil {
			log.Err(err).Msgf("unable to retrieve cluster name")
			return nil, err
		}
		err = json.Unmarshal([]byte(subnetsFromTofuOutput), &subnets)
		if err != nil {
			log.Err(err).Msgf("unable to retrieve subnets from tofu output")
			return nil, err
		}
		oidcProvider, err = tofu.GetOutput(ctx, k.clientConfig.ContextDirectory, "oidc_provider_arn", ".")
		if err != nil {
			log.Err(err).Msgf("unable to retrieve oidc provider from tofu output")
			return nil, err
		}
		oidcIssuer, err = tofu.GetOutput(ctx, k.clientConfig.ContextDirectory, "oidc_provider_url", ".")
		if err != nil {
			log.Err(err).Msgf("unable to retrieve oidc issuer from tofu output")
			return nil, err
		}
	}
	return map[string]interface{}{
		"cluster_name":      clusterName,
		"vpc":               vpc,
		"subnets":           subnets,
		"oidc_provider_arn": oidcProvider,
		"oidc_provider_url": oidcIssuer,
	}, nil
}

// PreExec implements ResourceClient.
func (k *kubernetesClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(k.clientConfig.Provider.Name(), k.clientConfig.ResourceName), k.clientConfig.ContextDirectory)
	if err != nil {
		return err
	}
	profile := k.configStore.GetString(ctx, store.ProfileKey)
	filesCopied, err := profiles.CopyInstallerFiles(profile, k.clientConfig.ContextDirectory)
	if err != nil {
		log.Err(err).Msgf("error copying profile files")
		return err
	}
	for _, file := range filesCopied {
		profileConfig := make(map[string]interface{})
		configData, err := os.ReadFile(file)
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
		if err != nil {
			log.Err(err).Msgf("unable to add profile config to datastore")
		}
	}
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
