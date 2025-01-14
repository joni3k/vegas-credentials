package cmd

import (
	"log"

	"github.com/aripalo/vegas-credentials/internal/commands/assume"
	"github.com/aripalo/vegas-credentials/internal/config"
	"github.com/spf13/cobra"
)

func init() {

	app, err := assume.New()
	if err != nil {
		log.Fatal(err)
	}
	cmd := buildAssumeCommand(app)

	rootCmd.AddCommand(cmd)

}

func buildAssumeCommand(app *assume.App) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "assume",
		Short: config.ASSUME_DESCRIPTION_SHORT,
		Long:  config.ASSUME_DESCRIPTION_LONG,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return app.PreRunE(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			app.Run()
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return app.PostRunE()
		},
	}

	cmd.Flags().String(
		config.Defaults.Profile.Name,
		config.Defaults.Profile.Value,
		config.Defaults.Profile.Usage,
	)
	err := cmd.MarkFlagRequired(config.Defaults.Profile.Name)
	if err != nil {
		panic(err)
	}

	cmd.Flags().Bool(
		config.Defaults.HideArns.Name,
		config.Defaults.HideArns.Value,
		config.Defaults.HideArns.Usage,
	)
	cmd.Flags().Bool(
		config.Defaults.DisableDialog.Name,
		config.Defaults.DisableDialog.Value,
		config.Defaults.DisableDialog.Usage,
	)
	cmd.Flags().Bool(
		config.Defaults.DisableMandatoryRefresh.Name,
		config.Defaults.DisableMandatoryRefresh.Value,
		config.Defaults.DisableMandatoryRefresh.Usage,
	)

	return cmd

}
