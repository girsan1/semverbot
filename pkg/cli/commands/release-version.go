package commands

import (
	"github.com/restechnica/semverbot/pkg/core"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/restechnica/semverbot/pkg/cli"
)

// NewReleaseVersionCommand creates a new release version command.
// Returns the new spf13/cobra command.
func NewReleaseVersionCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:     "version",
		PreRunE: ReleaseVersionCommandPreRunE,
		RunE:    ReleaseVersionCommandRunE,
	}

	command.Flags().StringVarP(&cli.ModeFlag, "mode", "m", "", "sbot mode")

	return command
}

// ReleaseVersionCommandPreRunE runs before the command runs.
// Returns an error if it fails.
func ReleaseVersionCommandPreRunE(cmd *cobra.Command, args []string) (err error) {
	return viper.BindPFlag(cli.SemverModeConfigKey, cmd.Flags().Lookup("mode"))
}

// ReleaseVersionCommandRunE runs the command.
// Returns an error if the command fails.
func ReleaseVersionCommandRunE(cmd *cobra.Command, args []string) error {
	var options = &core.ReleaseVersionOptions{
		DefaultVersion: cli.DefaultVersion,
		GitTagsPrefix:  viper.GetString(cli.GitTagsPrefixConfigKey),
		SemverMatchMap: viper.GetStringMapStringSlice(cli.SemverMatchConfigKey),
		SemverMode:     viper.GetString(cli.SemverModeConfigKey),
	}

	return core.ReleaseVersion(options)
}
