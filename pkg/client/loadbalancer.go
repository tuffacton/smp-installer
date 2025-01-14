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
	resourceName string
	configStore  store.DataStore
	outputStore  store.DataStore
}

func NewLoadbalancerClient(resourceName string, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &loadbalancerClient{
		resourceName: resourceName,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}

// Exec implements ResourceClient.
func (l *loadbalancerClient) Exec(ctx context.Context) error {
	isManaged, err := CheckIsManaged(ctx, l.resourceName, l.configStore)
	if err != nil {
		return err
	}
	if !isManaged {
		log.Info().Msgf("skipping %s sync as its not set to managed", l.resourceName)
		return nil
	}
	outDir := l.configStore.GetString(ctx, store.OutputDirectoryKey)
	provider := l.configStore.GetString(ctx, store.ProviderKey)
	contextDir := path.Join(outDir, provider, l.resourceName)
	tofu.ExecuteCommand(ctx, contextDir, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, contextDir, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (l *loadbalancerClient) PostExec(ctx context.Context) error {
	isManaged, err := CheckIsManaged(ctx, l.resourceName, l.configStore)
	if err != nil {
		return err
	}
	if !isManaged {
		existingIp := l.configStore.GetString(ctx, "loadbalancer.ip")
		l.outputStore.Set(ctx, "loadbalancer.ip", existingIp)
		return nil
	}
	outDir := l.configStore.GetString(ctx, store.OutputDirectoryKey)
	provider := l.configStore.GetString(ctx, store.ProviderKey)
	contextDir := path.Join(outDir, provider, l.resourceName)
	lbip, err := tofu.GetOutput(ctx, contextDir, "lb_eip_public_ip", ".")
	if err != nil {
		log.Err(err).Msgf("unable to retrieve loadbalancer ip")
		return err
	}
	l.outputStore.Set(ctx, "loadbalancer.ip", lbip)
	return nil
}

// PreExec implements ResourceClient.
func (l *loadbalancerClient) PreExec(ctx context.Context) error {
	outDir := l.configStore.GetString(ctx, store.OutputDirectoryKey)
	provider := l.configStore.GetString(ctx, store.ProviderKey)
	contextDir := path.Join(outDir, provider, l.resourceName)
	err := tofu.CopyFiles(path.Join(provider, l.resourceName), contextDir)
	if err != nil {
		return err
	}
	renderer := render.NewTemplateRenderer(l.configStore, l.outputStore)
	return renderer.Render(ctx, make(map[string]interface{}), "Resource", contextDir)
}
