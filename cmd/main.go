package main

import (
	"context"
	"os"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/cmd/resources"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func main() {
	cmd := &cobra.Command{
		Use:          "harness-smp",
		Short:        "Harness SMP Installer.",
		Long:         "Harness SMP Installer.",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringP("config", "c", "", "configuration file path")
	cmd.PersistentFlags().Bool("dry-run", false, "output only tf files and doesn't execute them")

	cmd.AddCommand(newSyncCommand())

	if err := cmd.Execute(); err != nil {
		log.Err(err).Msg("command failed")
		os.Exit(1)
	}
}

func newSyncCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sync",
		Short:        "Perform infrastructure operations using OpenTofu",
		Long:         "Perform infrastructure operations using OpenTofu",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			config := make(map[string]interface{})
			configFilePath, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}
			buf, err := os.ReadFile(configFilePath)
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(buf, config)
			if err != nil {
				return err
			}
			ctx := context.Background()
			configStore := store.NewMemoryStoreWithData(config)
			outputStore := store.NewMemoryStore()
			allCmds := getResourceCommands()
			for _, cmd := range allCmds {
				err = cmd.Sync(ctx, configStore, outputStore)
				if err != nil {
					log.Err(err).Msgf("step %s failed", cmd.Name())
					return err
				}
			}
			return nil
		},
	}
	return cmd
}

func getResourceCommands() []resources.ResourceCommand {
	resourceCommands := make([]resources.ResourceCommand, 0)
	resourceCommands = append(resourceCommands, resources.NewLoadbalancerCommand())
	resourceCommands = append(resourceCommands, resources.NewKubernetesCommand())
	resourceCommands = append(resourceCommands, resources.NewHarnessCommand())
	return resourceCommands
}

// func syncInfra(config *client.Configuration, cmd *cobra.Command) error {
// 	config.InfraConfig.OutputDir = path.Join(config.OutputDirectory, "infra")
// 	infraClient := client.NewKubernetesClient(&config.InfraConfig)
// 	err := infraClient.CopyTofuFiles()
// 	if err != nil {
// 		return err
// 	}
// 	err = infraClient.PrepareInputVariables()
// 	if err != nil {
// 		return err
// 	}
// 	infraClient.ExecuteInfraCommand(client.InitCommand)
// 	dryRun, _ := cmd.Flags().GetBool("dry-run")
// 	if dryRun {
// 		return infraClient.ExecuteInfraCommand(client.PlanCommand)
// 	}
// 	return infraClient.ExecuteInfraCommand(client.ApplyCommand)
// }

// func populateHarnessConfig(config *client.Configuration, cmd *cobra.Command) error {
// 	config.InfraConfig.OutputDir = path.Join(config.OutputDirectory, "infra")
// 	config.HarnessConfig.OutputDirectory = path.Join(config.OutputDirectory, "helm")
// 	infraClient := client.NewKubernetesClient(&config.InfraConfig)
// 	infraout, err := infraClient.TofuOutput()
// 	log.Info().Msgf("k8s endpoint: %s", infraout.KubernetesClusterEndpoint)
// 	if err != nil {
// 		return err
// 	}
// 	config.HarnessConfig.KubernetesClusterCertificate = infraout.KubernetesClusterCertificate
// 	config.HarnessConfig.KubernetesClusterEndpoint = infraout.KubernetesClusterEndpoint
// 	config.HarnessConfig.KubernetesClusterToken = infraout.KubernetesClusterToken
// 	config.HarnessConfig.LoadbalancerIP = infraout.LoadbalancerIP
// 	return nil
// }

// func syncHarness(config *client.Configuration, cmd *cobra.Command) error {
// 	config.HarnessConfig.OutputDirectory = path.Join(config.OutputDirectory, "helm")
// 	helmClient := client.NewHelmClient(&config.HarnessConfig)
// 	err := helmClient.CopyTofuFiles()
// 	if err != nil {
// 		return err
// 	}
// 	err = helmClient.PrepareOverride()
// 	if err != nil {
// 		return err
// 	}
// 	err = helmClient.PrepareInputVariables()
// 	if err != nil {
// 		return err
// 	}
// 	dryRun, _ := cmd.Flags().GetBool("dry-run")
// 	helmClient.Install(client.InitCommand)
// 	if dryRun {
// 		return helmClient.Install(client.PlanCommand)
// 	}
// 	return helmClient.Install(client.ApplyCommand)
// }
