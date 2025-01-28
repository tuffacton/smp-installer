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

type helmClient struct {
	clientConfig ClientConfig
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (h *helmClient) Exec(ctx context.Context) error {
	if !h.clientConfig.IsManaged {
		log.Info().Msgf("skipping %s sync as its not set to managed", h.clientConfig.ResourceName)
		return nil
	}
	tofu.ExecuteCommand(ctx, h.clientConfig.ContextDirectory, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, h.clientConfig.ContextDirectory, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (h *helmClient) PostExec(ctx context.Context) (map[string]interface{}, error) {
	return nil, nil
}

// PreExec implements ResourceClient.
func (h *helmClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(h.clientConfig.Provider.Name(), h.clientConfig.ResourceName), h.clientConfig.ContextDirectory)
	if err != nil {
		log.Err(err).Msgf("unable to copy files for helm chart module")
		return err
	}
	profile := h.configStore.GetString(ctx, store.ProfileKey)
	err = profiles.CopyFiles(profile, h.clientConfig.ContextDirectory, []string{"override.yaml"})
	if err != nil {
		log.Err(err).Msgf("unable to copy files for helm chart module")
		return err
	}
	resourceData := make(map[string]interface{})
	resourceData["override_file"] = "override.yaml"
	renderer := render.NewTemplateRenderer(h.configStore, h.outputStore)
	err = renderer.Render(ctx, resourceData, h.clientConfig.ContextDirectory)
	if err != nil {
		log.Err(err).Msgf("could not prepare override file")
		return err
	}
	return h.prepareOverride(ctx, path.Join(h.clientConfig.ContextDirectory, "override.yaml"))
}

func (h *helmClient) prepareOverride(ctx context.Context, profileYamlPath string) error {
	overrides, err := h.mergeOverrides(ctx, profileYamlPath)
	if err != nil {
		log.Err(err).Msgf("unable to merge existing override files")
		return err
	}
	data, err := yaml.Marshal(overrides)
	if err != nil {
		log.Err(err).Msgf("unable to serialize override configuration")
		return err
	}
	return os.WriteFile(profileYamlPath, data, 0666)
}

func (h *helmClient) mergeOverrides(ctx context.Context, profileYamlPath string) (map[string]interface{}, error) {
	value, err := h.configStore.Get(ctx, "harness.override_files")
	overrideFiles := make([]string, 0)
	overrideFiles = append(overrideFiles, profileYamlPath)
	if err != nil {
		log.Err(err).Msgf("unable to retrieve existing override files")
	} else {
		existingOverrides, ok := value.([]string)
		if ok {
			overrideFiles = append(overrideFiles, existingOverrides...)
		}
	}
	var master = make(map[string]interface{})
	for _, f := range overrideFiles {
		bs, err := os.ReadFile(f)
		var override = make(map[string]interface{})
		if err != nil {
			log.Err(err).Msgf("could not read file %s", f)
			return nil, err
		}
		if err := yaml.Unmarshal(bs, &override); err != nil {
			log.Err(err).Msgf("could not read yaml %s", f)
			return nil, err
		}

		for k, v := range override {
			master[k] = v
		}
	}
	return master, nil
}

func NewHelmClient(clientConfig ClientConfig, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &helmClient{
		clientConfig: clientConfig,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}
