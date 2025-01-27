package client

import (
	"context"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/profiles"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"github.com/rs/zerolog/log"
)

type dnsClient struct {
	configStore  store.DataStore
	outputStore  store.DataStore
	clientConfig ClientConfig
}

// Exec implements ResourceClient.
func (k *dnsClient) Exec(ctx context.Context) error {
	if !k.clientConfig.IsManaged {
		log.Info().Msgf("skipping %s sync as its not set to managed", k.clientConfig.ResourceName)
		return nil
	}
	tofu.ExecuteCommand(ctx, k.clientConfig.ContextDirectory, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, k.clientConfig.ContextDirectory, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (k *dnsClient) PostExec(ctx context.Context) (map[string]interface{}, error) {
	domainName := k.configStore.GetString(ctx, "dns.domain")
	if len(domainName) == 0 {
		domainName = k.outputStore.GetString(ctx, "loadbalancer.ip")
	}
	return map[string]interface{}{"domain": domainName}, nil
}

// PreExec implements ResourceClient.
func (k *dnsClient) PreExec(ctx context.Context) error {
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
	renderer := render.NewTemplateRenderer(k.configStore, k.outputStore)
	return renderer.Render(ctx, map[string]interface{}{}, k.clientConfig.ContextDirectory)
}

func NewDnsClient(conf ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &dnsClient{
		clientConfig: conf,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}
