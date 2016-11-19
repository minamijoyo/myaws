package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws/rds"
)

func init() {
	RootCmd.AddCommand(newRDSCmd())
}

func newRDSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rds",
		Short: "Manage RDS resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newRDSLsCmd(),
	)

	return cmd
}

func newRDSLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List RDS instances",
		RunE:  rds.Ls,
	}

	flags := cmd.Flags()
	flags.BoolP("quiet", "q", false, "Only display DBInstanceIdentifier")
	flags.StringP("fields", "F", "DBInstanceClass Engine AllocatedStorage StorageTypeIops InstanceCreateTime DBInstanceIdentifier ReadReplicaSource", "Output fields list separated by space")

	viper.BindPFlag("rds.ls.quiet", flags.Lookup("quiet"))
	viper.BindPFlag("rds.ls.fields", flags.Lookup("fields"))

	return cmd
}
