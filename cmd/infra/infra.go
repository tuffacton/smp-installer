package infra

import (
	"io/ioutil"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewInfraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "infra plan|apply",
		Short:        "Perform infrastructure operations using OpenTofu",
		Long:         "Perform infrastructure operations using OpenTofu",
		SilenceUsage: true,
	}
	cmd.AddCommand(newPlanCommand())
	cmd.AddCommand(newApplyCommand())
	return cmd
}

func newPlanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "plan",
		Short:        "Plan the infra resources using OpenTofu",
		Long:         "Plan the infra resources using OpenTofu",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := new(client.Configuration)
			configFilePath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}
			buf, err := ioutil.ReadFile(configFilePath)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(buf, config)
			if err != nil {
				return err
			}
			infraClient := client.NewInfraClient(&config.InfraConfig)
			err = infraClient.CopyTofuFiles()
			if err != nil {
				return err
			}
			err = infraClient.PrepareInputVariables()
			if err != nil {
				return err
			}
			err = infraClient.ExecuteInfraCommand(client.InitCommand)
			if err != nil {
				return err
			}
			return infraClient.ExecuteInfraCommand(client.PlanCommand)
		},
	}

	return cmd
}

func newApplyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "apply",
		Short:        "Apply the infra resources using OpenTofu",
		Long:         "Apply the infra resources using OpenTofu",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := new(client.Configuration)
			configFilePath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}
			buf, err := ioutil.ReadFile(configFilePath)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(buf, config)
			if err != nil {
				return err
			}
			infraClient := client.NewInfraClient(&config.InfraConfig)
			return infraClient.ExecuteInfraCommand(client.ApplyCommand)
		},
	}

	return cmd
}
