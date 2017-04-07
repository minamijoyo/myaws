package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws"
)

var cfgFile string

// RootCmd is a top level command instance
var RootCmd = &cobra.Command{
	Use:           "myaws",
	Short:         "A human friendly AWS CLI written in Go.",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.myaws.yml)")
	RootCmd.PersistentFlags().BoolP("debug", "", false, "Enable debug mode")
	RootCmd.PersistentFlags().StringP("profile", "", "", "AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)")
	RootCmd.PersistentFlags().StringP("region", "", "", "AWS region (default none and used AWS_DEFAULT_REGION environment variable.")
	RootCmd.PersistentFlags().StringP("timezone", "", "Local", "Time zone, such as UTC, Asia/Tokyo")
	RootCmd.PersistentFlags().BoolP("humanize", "", true, "Use Human friendly format for time")
	RootCmd.PersistentFlags().StringP("format", "", "tsv", "Output format")

	viper.BindPFlag("profile", RootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("region", RootCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("timezone", RootCmd.PersistentFlags().Lookup("timezone"))
	viper.BindPFlag("humanize", RootCmd.PersistentFlags().Lookup("humanize"))
	viper.BindPFlag("format", RootCmd.PersistentFlags().Lookup("format"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".myaws")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

func newClient() (*myaws.Client, error) {
	return myaws.NewClient(
		os.Stdin,
		os.Stdout,
		os.Stderr,
		viper.GetString("profile"),
		viper.GetString("region"),
		viper.GetString("timezone"),
		viper.GetBool("humanize"),
		viper.GetString("format"),
	)
}
