package client

import (
	"context"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"github.com/rs/zerolog/log"
)

type loadbalancerClient struct {
	clientConfig ClientConfig
	configStore  store.DataStore
	outputStore  store.DataStore
}

func NewLoadbalancerClient(clientConfig ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &loadbalancerClient{
		clientConfig: clientConfig,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}

// Exec implements ResourceClient.
func (l *loadbalancerClient) Exec(ctx context.Context) error {
	if !l.clientConfig.IsManaged {
		log.Info().Msgf("skipping %s sync as its not set to managed", l.clientConfig.ResourceName)
		return nil
	}
	tofu.ExecuteCommand(ctx, l.clientConfig.ContextDirectory, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, l.clientConfig.ContextDirectory, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (l *loadbalancerClient) PostExec(ctx context.Context) (map[string]interface{}, error) {
	if !l.clientConfig.IsManaged {
		existingIp := l.configStore.GetString(ctx, "loadbalancer.ip")
		l.outputStore.Set(ctx, "loadbalancer.ip", existingIp)
		return nil, nil
	}
	lbip, err := tofu.GetOutput(ctx, l.clientConfig.ContextDirectory, "dns_name", ".")
	if err != nil {
		log.Err(err).Msgf("unable to retrieve loadbalancer ip")
		return nil, err
	}
	lbzone, err := tofu.GetOutput(ctx, l.clientConfig.ContextDirectory, "dns_zone_id", ".")
	if err != nil {
		log.Err(err).Msgf("unable to retrieve loadbalancer zone")
		return nil, err
	}
	return map[string]interface{}{"ip": lbip, "zone_id": lbzone}, nil
}

// PreExec implements ResourceClient.
func (l *loadbalancerClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(l.clientConfig.Provider.Name(), l.clientConfig.ResourceName), l.clientConfig.ContextDirectory)
	if err != nil {
		return err
	}
	renderer := render.NewTemplateRenderer(l.configStore, l.outputStore)
	return renderer.Render(ctx, make(map[string]interface{}), l.clientConfig.ContextDirectory)
}
