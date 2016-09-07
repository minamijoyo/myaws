package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "myaws",
	Short: "myaws is a simple command line tool for managing my aws resources",
}

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

	viper.BindPFlag("profile", RootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("region", RootCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("timezone", RootCmd.PersistentFlags().Lookup("timezone"))
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
