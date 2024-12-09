package main

import (
	"os"

	"git0.harness.io/l7B_kbSEQD2wjrM7PShm5w/PROD/Harness_Commons/harness-smp-installer/cmd/infra"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {

	cmd := &cobra.Command{
		Use:          "harness-smp",
		Short:        "Harness SMP Installer.",
		Long:         "Harness SMP Installer.",
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringP("config", "c", "", "configuration file path")

	cmd.AddCommand(infra.NewInfraCommand())

	if err := cmd.Execute(); err != nil {
		log.Err(err).Msg("command failed")
		os.Exit(1)
	}

}
