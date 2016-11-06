package cmd

import (
	"github.com/minamijoyo/myaws/myaws/rds"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rdsLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List RDS instances",
	Run:   rds.Ls,
}

func init() {
	rdsCmd.AddCommand(rdsLsCmd)

	rdsLsCmd.Flags().BoolP("quiet", "q", false, "Only display DBInstanceIdentifier")
	rdsLsCmd.Flags().StringP("fields", "F", "DBInstanceClass Engine AllocatedStorage StorageTypeIops InstanceCreateTime DBInstanceIdentifier ReadReplicaSource", "Output fields list separated by space")

	viper.BindPFlag("rds.ls.quiet", rdsLsCmd.Flags().Lookup("quiet"))
	viper.BindPFlag("rds.ls.fields", rdsLsCmd.Flags().Lookup("fields"))
}
