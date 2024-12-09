package client

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"text/template"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/tofu"
	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
)

type InfraCommand string

const (
	InitCommand  InfraCommand = "init"
	PlanCommand  InfraCommand = "plan"
	ApplyCommand InfraCommand = "apply"
)

type InfraConfiguration struct {
	TofuPlanOutputDir string                  `yaml:"plan_output_dir"`
	TofuBinaryPath    string                  `yaml:"tofu_path"`
	CloudProvider     CloudProviderType       `yaml:"provider"`
	KubernetesConfig  KubernetesConfiguration `yaml:"kubernetes"`
	NodesConfig       NodesConfig             `yaml:"nodes"`
	Region            string                  `yaml:"region"`
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

type InfraClient struct {
	Config *InfraConfiguration
}

func NewInfraClient(conf *InfraConfiguration) *InfraClient {
	return &InfraClient{Config: conf}
}

func (i *InfraClient) CopyTofuFiles() error {
	if i.Config.CloudProvider == AWSCloudProvider {
		os.Mkdir(i.Config.TofuPlanOutputDir, 0777)
		files, err := tofu.AWSTofuFiles.ReadDir("aws/infra")
		if err != nil {
			log.Err(err).Msg("could not find infra directory")
			return err
		}
		for _, f := range files {
			data, err := tofu.AWSTofuFiles.ReadFile(path.Join("aws/infra", f.Name()))
			if err != nil {
				log.Err(err).Msgf("cannot read file %s", f.Name())
				return err
			}
			err = os.WriteFile(path.Join(i.Config.TofuPlanOutputDir, f.Name()), data, 0666)
			if err != nil {
				log.Err(err).Msgf("failed to copy file %s to directory %s", f.Name(), i.Config.TofuPlanOutputDir)
				return err
			}
		}
	}
	return nil
}

func (i *InfraClient) PrepareInputVariables() error {
	if i.Config.CloudProvider == AWSCloudProvider {
		variableFilePath := path.Join(i.Config.TofuPlanOutputDir, "tf.vars.tpl")
		tmpl := template.New("tf.vars.tpl")
		variablesFile, err := os.Create(path.Join(i.Config.TofuPlanOutputDir, "tf.vars"))
		if err != nil {
			log.Err(err).Msg("could not create tofu variables file")
			return err
		}
		tmpl.Funcs(sprig.FuncMap())
		tmpl.ParseFiles(variableFilePath)
		err = tmpl.ExecuteTemplate(variablesFile, "tf.vars.tpl", i.Config)
		if err != nil {
			// fd, err := os.ReadFile(variableFilePath)
			// if err != nil {
			// 	log.Err(err).Msgf("could not read variable template file: %s", variableFilePath)
			// } else {
			// 	log.Info().Msgf("variable template file: %s", fd)
			// }
			log.Err(err).Msgf("could not execute template to format variables")
			return err
		}
	}
	return nil
}

func (i *InfraClient) ExecuteInfraCommand(icmd InfraCommand) error {
	os.Chdir(i.Config.TofuPlanOutputDir)
	cmd := exec.Command("tofu", fmt.Sprintf("%s", icmd), "-var-file=tf.vars", "-auto-approve")
	stdout, err := cmd.Output()

	if err != nil {
		log.Err(err).Msgf("%s command failed with output: %s", icmd, err.(*exec.ExitError).Stderr)
		return err
	}

	log.Info().Msgf("%s command output", icmd)
	log.Info().Msgf("%s", stdout)

	return nil
}
