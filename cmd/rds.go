package cmd

import "github.com/spf13/cobra"

var rdsCmd = &cobra.Command{
	Use:   "rds",
	Short: "Manage RDS resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(rdsCmd)
}
