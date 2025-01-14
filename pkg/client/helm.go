package client

import (
	"context"
	"os"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/profiles"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type HarnessConfiguration struct {
	Version                   string   `yaml:"version"`
	ExistingOverrideFilePaths []string `yaml:"override_files"`
	Namespace                 string   `yaml:"namespace"`

	OutputDirectory              string
	KubernetesClusterCertificate string
	KubernetesClusterEndpoint    string
	KubernetesClusterToken       string
	LoadbalancerIP               string
}

type HelmTerraformConfiguration struct {
	OverrideFile                 string
	Registry                     string
	KubernetesClusterCertificate string
	KubernetesClusterEndpoint    string
	KubernetesClusterToken       string
	Namespace                    string
	Version                      string
}

type helmClient struct {
	resourceName string
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (h *helmClient) Exec(ctx context.Context) error {
	outDir := h.configStore.GetString(ctx, store.OutputDirectoryKey)
	provider := h.configStore.GetString(ctx, store.ProviderKey)
	contextDir := path.Join(outDir, provider, h.resourceName)
	tofu.ExecuteCommand(ctx, contextDir, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, contextDir, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (h *helmClient) PostExec(ctx context.Context) error {
	return nil
}

// PreExec implements ResourceClient.
func (h *helmClient) PreExec(ctx context.Context) error {
	outDir := h.configStore.GetString(ctx, store.OutputDirectoryKey)
	provider := h.configStore.GetString(ctx, store.ProviderKey)
	contextDir := path.Join(outDir, provider, h.resourceName)
	err := tofu.CopyFiles(path.Join(provider, h.resourceName), contextDir)
	if err != nil {
		log.Err(err).Msgf("unable to copy files for helm chart module")
		return err
	}
	profile := h.configStore.GetString(ctx, store.ProfileKey)
	err = profiles.CopyFiles(profile, contextDir)
	if err != nil {
		log.Err(err).Msgf("unable to copy files for helm chart module")
		return err
	}
	resourceData := make(map[string]interface{})
	resourceData["override_file"] = "override-" + profile + ".yaml"
	renderer := render.NewTemplateRenderer(h.configStore, h.outputStore)
	err = renderer.Render(ctx, resourceData, "Resource", contextDir)
	if err != nil {
		log.Err(err).Msgf("could not prepare override file")
		return err
	}
	return h.prepareOverride(ctx, path.Join(contextDir, "override-"+profile+".yaml"))
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

func NewHelmClient(resourceName string, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &helmClient{
		resourceName: resourceName,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}

// func (h *HelmClient) CopyTofuFiles() error {
// 	err := os.MkdirAll(h.config.OutputDirectory, 0777)
// 	if err != nil {
// 		log.Err(err).Msgf("could not create output directory: %s", h.config.OutputDirectory)
// 		return err
// 	}
// 	files, err := tofu.AWSTofuFiles.ReadDir("aws/helm")
// 	if err != nil {
// 		log.Err(err).Msg("could not find helm directory")
// 		return err
// 	}
// 	for _, f := range files {
// 		data, err := tofu.AWSTofuFiles.ReadFile(path.Join("aws/helm", f.Name()))
// 		if err != nil {
// 			log.Err(err).Msgf("cannot read file %s", f.Name())
// 			return err
// 		}
// 		err = os.WriteFile(path.Join(h.config.OutputDirectory, f.Name()), data, 0666)
// 		if err != nil {
// 			log.Err(err).Msgf("failed to copy file %s to directory %s", f.Name(), h.config.OutputDirectory)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (h *HelmClient) PrepareInputVariables() error {
// 	// finalOverrideFilePath := path.Join(h.config.OutputDirectory, "harness-override.yaml")
// 	templateConfig := HelmTerraformConfiguration{
// 		OverrideFile:                 "harness-override.yaml",
// 		Namespace:                    h.config.Namespace,
// 		KubernetesClusterCertificate: h.config.KubernetesClusterCertificate,
// 		KubernetesClusterEndpoint:    h.config.KubernetesClusterEndpoint,
// 		KubernetesClusterToken:       h.config.KubernetesClusterToken,
// 		Registry:                     "https://harness.github.io/helm-charts",
// 		Version:                      h.config.Version,
// 	}
// 	variableFilePath := path.Join(h.config.OutputDirectory, "tf.vars.tpl")
// 	tmpl := template.New("tf.vars.tpl")
// 	variablesFile, err := os.Create(path.Join(h.config.OutputDirectory, "tf.vars"))
// 	if err != nil {
// 		log.Err(err).Msg("could not create tofu variables file")
// 		return err
// 	}
// 	tmpl.Funcs(sprig.FuncMap())
// 	tmpl.ParseFiles(variableFilePath)
// 	err = tmpl.ExecuteTemplate(variablesFile, "tf.vars.tpl", templateConfig)
// 	if err != nil {
// 		log.Err(err).Msgf("could not execute template to format variables")
// 		return err
// 	}
// 	return nil
// }

// func (h *HelmClient) PrepareOverride() error {
// 	overrides, err := h.mergeOverrides()
// 	if err != nil {
// 		log.Err(err).Msgf("unable to merge existing override files")
// 		return err
// 	}
// 	h.setValue(overrides, "platform.bootstrap.networking.nginx.loadBalancerIP", h.config.LoadbalancerIP)
// 	h.setValue(overrides, "platform.bootstrap.networking.nginx.loadBalancerEnabled", true)
// 	h.setValue(overrides, "global.loadbalancerURL", fmt.Sprintf("http://%s", h.config.LoadbalancerIP))
// 	finalOverrideFilePath := path.Join(h.config.OutputDirectory, "harness-override.yaml")
// 	data, err := yaml.Marshal(overrides)
// 	if err != nil {
// 		log.Err(err).Msgf("unable to serialize override configuration")
// 		return err
// 	}
// 	return os.WriteFile(finalOverrideFilePath, data, 0666)
// }

// func (h *HelmClient) setValue(overrides map[string]interface{}, path string, value any) error {
// 	segments := strings.Split(path, ".")
// 	currentOverrides := overrides
// 	for idx, seg := range segments {
// 		// log.Info().Msgf("checking for segment: %s", seg)
// 		if idx == len(segments)-1 {
// 			currentOverrides[seg] = value
// 			return nil
// 		}
// 		_, ok := currentOverrides[seg]
// 		if !ok {
// 			currentOverrides[seg] = make(map[string]interface{})
// 		}
// 		currentOverrides = currentOverrides[seg].(map[string]interface{})
// 	}
// 	return nil
// }

// func (h *HelmClient) Install(icmd InfraCommand) error {
// 	currentDir, _ := os.Getwd()
// 	os.Chdir(h.config.OutputDirectory)
// 	extraArgs := []string{"-auto-approve"}
// 	cmd := exec.Command("tofu", string(icmd), "-var-file=tf.vars", "-no-color")
// 	if icmd != InitCommand && icmd != PlanCommand {
// 		cmd.Args = append(cmd.Args, extraArgs...)
// 	}
// 	stdout, err := cmd.Output()
// 	if err != nil {
// 		log.Err(err).Msgf("%s command failed with output: %s", icmd, err.(*exec.ExitError).Stderr)
// 		return err
// 	}

// 	log.Info().Msgf("%s command output", icmd)
// 	log.Info().Msgf("%s", string(stdout))
// 	os.Chdir(currentDir)
// 	return nil
// }
