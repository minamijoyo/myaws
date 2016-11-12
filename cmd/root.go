package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd is a top level command instance
var RootCmd = &cobra.Command{
	Use:   "myaws",
	Short: "A human friendly aws cli written in Go.",
}

// Execute is a entrypoint of cli
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default $HOME/.myaws.yml)")
	RootCmd.PersistentFlags().StringP("profile", "", "", "AWS profile (default none and used AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY environment variables.)")
	RootCmd.PersistentFlags().StringP("region", "", "", "AWS region (default none and used AWS_DEFAULT_REGION environment variable.")
	RootCmd.PersistentFlags().StringP("timezone", "", "Local", "Time zone, such as UTC, Asia/Tokyo")
	RootCmd.PersistentFlags().BoolP("humanize", "", true, "Use Human friendly format for time")

	viper.BindPFlag("profile", RootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("region", RootCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("timezone", RootCmd.PersistentFlags().Lookup("timezone"))
	viper.BindPFlag("humanize", RootCmd.PersistentFlags().Lookup("humanize"))

	RootCmd.AddCommand(newEC2Cmd())
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
