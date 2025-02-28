package client

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/harness/smp-installer/pkg/profiles"
	"github.com/harness/smp-installer/pkg/render"
	"github.com/harness/smp-installer/pkg/store"
	"github.com/harness/smp-installer/pkg/tofu"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type vpcClient struct {
	clientConfig ClientConfig
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (v *vpcClient) Exec(ctx context.Context) error {
	if !v.clientConfig.IsManaged {
		log.Info().Msgf("skipping %s sync as its not set to managed", v.clientConfig.ResourceName)
		return nil
	}
	tofu.ExecuteCommand(ctx, v.clientConfig.ContextDirectory, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, v.clientConfig.ContextDirectory, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (v *vpcClient) PostExec(ctx context.Context) (output map[string]interface{}, err error) {
	vpc := ""
	subnets := []string{}
	privateSubnets := []string{}
	azs := []string{}
	if !v.clientConfig.IsManaged {
		vpc = v.configStore.GetString(ctx, "vpc.id")
		subnetsFromConfig, err := v.configStore.Get(ctx, "vpc.subnets.public")
		if err == nil {
			subnetsFromConfigArray := subnetsFromConfig.([]interface{})
			for _, subnet := range subnetsFromConfigArray {
				subnetObj := subnet.(map[string]interface{})
				subnets = append(subnets, subnetObj["id"].(string))
			}
		}
		privateSubnetsFromConfig, err := v.configStore.Get(ctx, "vpc.subnets.public")
		if err == nil {
			subnetsFromConfigArray := privateSubnetsFromConfig.([]interface{})
			for _, subnet := range subnetsFromConfigArray {
				subnetObj := subnet.(map[string]interface{})
				privateSubnets = append(privateSubnets, subnetObj["id"].(string))
			}
		}
		azsFromConfig, err := v.configStore.Get(ctx, "vpc.availability_zones")
		if err == nil {
			azsFromConfigArray := azsFromConfig.([]interface{})
			for _, az := range azsFromConfigArray {
				azs = append(azs, az.(string))
			}
		}
	} else {
		var err error = nil
		vpc, err = tofu.GetOutput(ctx, v.clientConfig.ContextDirectory, "vpc", ".")
		if err != nil {
			return nil, err
		}
		subnetsFromTofuOutput, err := tofu.GetOutput(ctx, v.clientConfig.ContextDirectory, "subnets", ".")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(subnetsFromTofuOutput), &subnets)
		if err != nil {
			log.Err(err).Msgf("unable to unmarshal subnets from tofu output")
			return nil, err
		}
		privateSubnetsFromTofuOutput, err := tofu.GetOutput(ctx, v.clientConfig.ContextDirectory, "private_subnets", ".")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(privateSubnetsFromTofuOutput), &privateSubnets)
		if err != nil {
			log.Err(err).Msgf("unable to unmarshal private_subnets from tofu output")
			return nil, err
		}
		azsFromTofuOutput, err := tofu.GetOutput(ctx, v.clientConfig.ContextDirectory, "availability_zones", ".")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(azsFromTofuOutput), &azs)
		if err != nil {
			log.Err(err).Msgf("unable to unmarshal availability_zones from tofu output")
			return nil, err
		}
	}
	return map[string]interface{}{
		"vpc":                vpc,
		"subnets":            subnets,
		"private_subnets":    privateSubnets,
		"availability_zones": azs,
	}, nil
}

// PreExec implements ResourceClient.
func (v *vpcClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(v.clientConfig.Provider.Name(), v.clientConfig.ResourceName), v.clientConfig.ContextDirectory)
	if err != nil {
		return err
	}
	profile := v.configStore.GetString(ctx, store.ProfileKey)
	filesCopied, err := profiles.CopyInstallerFiles(profile, v.clientConfig.ContextDirectory)
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
		err = v.configStore.AddAll(ctx, profileConfig)
		if err != nil {
			log.Err(err).Msgf("unable to add profile config to datastore")
		}
	}
	renderer := render.NewTemplateRenderer(v.configStore, v.outputStore)
	return renderer.Render(ctx, map[string]interface{}{}, v.clientConfig.ContextDirectory)
}

func NewVPCClient(clientConfig ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &vpcClient{
		clientConfig: clientConfig,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}
