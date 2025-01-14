package client

import (
	"context"
	"path"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/render"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"github.com/rs/zerolog/log"
)

type InfraCommand string

const (
	InitCommand  InfraCommand = "init"
	PlanCommand  InfraCommand = "plan"
	ApplyCommand InfraCommand = "apply"
)

type InfraConfiguration struct {
	OutputDir        string
	TofuBinaryPath   string                  `yaml:"tofu_path"`
	CloudProvider    CloudProviderType       `yaml:"provider"`
	KubernetesConfig KubernetesConfiguration `yaml:"kubernetes"`
	NodesConfig      NodesConfig             `yaml:"nodes"`
	Region           string                  `yaml:"region"`
}

type NodesConfig struct {
	MachineType  string `yaml:"machine_type"`
	MinimumNodes int    `yaml:"minimum_nodes"`
	MaximumNodes int    `yaml:"maximum_nodes"`
	InitialNodes int    `yaml:"initial_nodes"`
}

type KubernetesConfiguration struct {
	Version string `yaml:"version"`
}

type InfraOutput struct {
	KubernetesClusterCertificate string
	KubernetesClusterEndpoint    string
	KubernetesClusterToken       string
	LoadbalancerIP               string
}

type kubernetesClient struct {
	resourceName string
	configStore  store.DataStore
	outputStore  store.DataStore
}

// Exec implements ResourceClient.
func (k *kubernetesClient) Exec(ctx context.Context) error {
	isManaged, err := CheckIsManaged(ctx, k.resourceName, k.configStore)
	if err != nil {
		return err
	}
	if !isManaged {
		log.Info().Msgf("skipping %s sync as its not set to managed", k.resourceName)
		return nil
	}
	outDir := k.configStore.GetString(ctx, store.OutputDirectoryKey)
	provider := k.configStore.GetString(ctx, store.ProviderKey)
	contextDir := path.Join(outDir, provider, k.resourceName)
	tofu.ExecuteCommand(ctx, contextDir, tofu.InitCommand)
	return tofu.ExecuteCommand(ctx, contextDir, tofu.ApplyCommand)
}

// PostExec implements ResourceClient.
func (k *kubernetesClient) PostExec(ctx context.Context) error {
	isManaged, err := CheckIsManaged(ctx, k.resourceName, k.configStore)
	if err != nil {
		return err
	}
	if !isManaged {
		clusterName := k.configStore.GetString(ctx, "kubernetes.cluster_name")
		k.outputStore.Set(ctx, "kubernetes.cluster_name", clusterName)
		return nil
	}
	outDir := k.configStore.GetString(ctx, store.OutputDirectoryKey)
	provider := k.configStore.GetString(ctx, store.ProviderKey)
	contextDir := path.Join(outDir, provider, k.resourceName)
	clusterName, err := tofu.GetOutput(ctx, contextDir, "clustername", ".")
	if err != nil {
		log.Err(err).Msgf("unable to retrieve cluster name")
		return err
	}
	k.outputStore.Set(ctx, "kubernetes.cluster_name", clusterName)
	return nil
}

// PreExec implements ResourceClient.
func (k *kubernetesClient) PreExec(ctx context.Context) error {
	outDir := k.configStore.GetString(ctx, store.OutputDirectoryKey)
	provider := k.configStore.GetString(ctx, store.ProviderKey)
	contextDir := path.Join(outDir, provider, k.resourceName)
	err := tofu.CopyFiles(path.Join(provider, k.resourceName), contextDir)
	if err != nil {
		return err
	}
	renderer := render.NewTemplateRenderer(k.configStore, k.outputStore)
	return renderer.Render(ctx, make(map[string]interface{}), "Resource", contextDir)
}

func NewKubernetesClient(resourceName string, configStore store.DataStore, outputStore store.DataStore) ResourceClient {
	return &kubernetesClient{
		resourceName: resourceName,
		configStore:  configStore,
		outputStore:  outputStore,
	}
}

// func (i *InfraClient) CopyTofuFiles() error {
// 	if i.Config.CloudProvider == AWSCloudProvider {
// 		os.Mkdir(i.Config.OutputDir, 0777)
// 		files, err := tofu.AWSTofuFiles.ReadDir("aws/infra")
// 		if err != nil {
// 			log.Err(err).Msg("could not find infra directory")
// 			return err
// 		}
// 		for _, f := range files {
// 			data, err := tofu.AWSTofuFiles.ReadFile(path.Join("aws/infra", f.Name()))
// 			if err != nil {
// 				log.Err(err).Msgf("cannot read file %s", f.Name())
// 				return err
// 			}
// 			err = os.WriteFile(path.Join(i.Config.OutputDir, f.Name()), data, 0666)
// 			if err != nil {
// 				log.Err(err).Msgf("failed to copy file %s to directory %s", f.Name(), i.Config.OutputDir)
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// func (i *InfraClient) PrepareInputVariables() error {
// 	if i.Config.CloudProvider == AWSCloudProvider {
// 		variableFilePath := path.Join(i.Config.OutputDir, "tf.vars.tpl")
// 		tmpl := template.New("tf.vars.tpl")
// 		variablesFile, err := os.Create(path.Join(i.Config.OutputDir, "tf.vars"))
// 		if err != nil {
// 			log.Err(err).Msg("could not create tofu variables file")
// 			return err
// 		}
// 		tmpl.Funcs(sprig.FuncMap())
// 		tmpl.ParseFiles(variableFilePath)
// 		err = tmpl.ExecuteTemplate(variablesFile, "tf.vars.tpl", i.Config)
// 		if err != nil {
// 			log.Err(err).Msgf("could not execute template to format variables")
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (i *InfraClient) ExecuteInfraCommand(icmd InfraCommand) error {
// 	currentDir, _ := os.Getwd()
// 	os.Chdir(i.Config.OutputDir)
// 	extraArgs := []string{"-auto-approve"}
// 	cmd := exec.Command("tofu", fmt.Sprintf("%s", icmd), "-var-file=tf.vars", "-no-color")
// 	if icmd != InitCommand && icmd != PlanCommand {
// 		cmd.Args = append(cmd.Args, extraArgs...)
// 	}
// 	stdout, err := cmd.Output()
// 	if err != nil {
// 		log.Err(err).Msgf("infra %s command failed with output: %s", icmd, err.(*exec.ExitError).Stderr)
// 		return err
// 	}

// 	log.Info().Msgf("infra %s command output", icmd)
// 	log.Info().Msg(string(stdout))

// 	os.Chdir(currentDir)
// 	return nil
// }

// func (i *InfraClient) TofuOutput() (*InfraOutput, error) {
// 	currentDir, _ := os.Getwd()
// 	os.Chdir(i.Config.OutputDir)
// 	clcert, err := getTofuOutput("eksout", ".cluster_certificate_authority_data")
// 	if err != nil {
// 		log.Err(err).Msgf("could not get cluster_certificate_authority_data from tofu output")
// 		return nil, err
// 	}
// 	clendpoint, err := getTofuOutput("eksout", ".cluster_endpoint")
// 	if err != nil {
// 		log.Err(err).Msgf("could not get cluster_endpoint from tofu output")
// 		return nil, err
// 	}
// 	cltoken, err := getTofuOutput("authout", ".")
// 	if err != nil {
// 		log.Err(err).Msgf("could not get auth token from tofu output")
// 		return nil, err
// 	}
// 	loadbalancerIp, err := getTofuOutput("lb_eip_public_ip", ".")
// 	if err != nil {
// 		log.Err(err).Msgf("could not get load balancer ip from tofu output")
// 		return nil, err
// 	}
// 	os.Chdir(currentDir)
// 	return &InfraOutput{
// 		KubernetesClusterCertificate: clcert,
// 		KubernetesClusterEndpoint:    clendpoint,
// 		KubernetesClusterToken:       cltoken,
// 		LoadbalancerIP:               loadbalancerIp,
// 	}, nil
// }

// func getTofuOutput(outname, jspath string) (string, error) {
// 	tmpfilename := fmt.Sprintf("tmpout-%s.json", outname)
// 	_, err := os.Lstat(tmpfilename)
// 	if err != nil {
// 		tofucmd := exec.Command("tofu", "output", "-json", outname)
// 		tofuout, err := tofucmd.Output()
// 		if err != nil {
// 			log.Err(err).Msgf("tofu output failed")
// 		}
// 		log.Info().Msgf("Loadbalancer IP: %s", string(tofuout))
// 		err = os.WriteFile(tmpfilename, tofuout, 0666)
// 		if err != nil {
// 			log.Err(err).Msgf("failed to create tofu json output file")
// 			return "", nil
// 		}
// 	}
// 	jqcmd := exec.Command("jq", "-r", jspath, tmpfilename)
// 	jqout, err := jqcmd.Output()
// 	if err != nil {
// 		log.Err(err).Msgf("jq command wait failed: %s", err.(*exec.ExitError).Stderr)
// 		return "", err
// 	}
// 	return strings.TrimSpace(string(jqout)), nil
// }
