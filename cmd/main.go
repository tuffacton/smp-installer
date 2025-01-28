package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/cmd/resources"
	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/pkg/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func main() {
	setupLogging()
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

func setupLogging() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
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
	resourceCommands = append(resourceCommands, resources.NewKubernetesCommand())
	resourceCommands = append(resourceCommands, resources.NewLoadbalancerCommand())
	resourceCommands = append(resourceCommands, resources.NewDnsCommand())
	resourceCommands = append(resourceCommands, resources.NewHarnessCommand())
	return resourceCommands
}
