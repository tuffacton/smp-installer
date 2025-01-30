package client

import (
	"context"
	"os"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/profiles"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/util"
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
	if !h.clientConfig.IsManaged {
		return map[string]interface{}{}, nil
	}
	return map[string]interface{}{}, nil
	// manifestOutput, err := tofu.GetOutput(ctx, h.clientConfig.ContextDirectory, "manifests", ".")
	// if err != nil {
	// 	log.Err(err).Msgf("unable to retrieve helm release output")
	// 	return nil, err
	// }
	// manifestJson := make(map[string]interface{})
	// err = json.Unmarshal([]byte(manifestOutput), &manifestJson)
	// if err != nil {
	// 	log.Err(err).Msgf("unable to unmarshal helm release output")
	// 	return nil, err
	// }
	// return map[string]interface{}{"manifests": manifestJson}, nil
}

// PreExec implements ResourceClient.
func (h *helmClient) PreExec(ctx context.Context) error {
	err := tofu.CopyFiles(path.Join(h.clientConfig.Provider.Name(), h.clientConfig.ResourceName), h.clientConfig.ContextDirectory)
	if err != nil {
		log.Err(err).Msgf("unable to copy files for helm chart module")
		return err
	}
	profile := h.configStore.GetString(ctx, store.ProfileKey)
	filesCopied, err := profiles.CopyOverrideFiles(profile, h.clientConfig.ContextDirectory)
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
	return h.prepareOverride(ctx,
		path.Join(h.clientConfig.ContextDirectory, "override.yaml"),
		filesCopied)
}

func (h *helmClient) prepareOverride(ctx context.Context,
	finalYamlPath string,
	profileOverrideFiles []string) error {
	overrides, err := h.mergeOverrides(ctx, profileOverrideFiles)
	if err != nil {
		log.Err(err).Msgf("unable to merge existing override files")
		return err
	}
	data, err := yaml.Marshal(overrides)
	if err != nil {
		log.Err(err).Msgf("unable to serialize override configuration")
		return err
	}
	return os.WriteFile(finalYamlPath, data, 0666)
}

func (h *helmClient) mergeOverrides(ctx context.Context,
	profileOverrideFiles []string) (map[string]interface{}, error) {
	value, err := h.configStore.Get(ctx, "harness.override_files")
	overrideFiles := make([]string, 0)
	if err != nil {
		log.Err(err).Msgf("unable to retrieve existing override files")
	} else {
		existingOverrides, ok := value.([]string)
		if ok {
			overrideFiles = append(overrideFiles, existingOverrides...)
		}
	}
	overrideFiles = append(overrideFiles, profileOverrideFiles...)
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
		master = util.MergeMaps(master, override)
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
