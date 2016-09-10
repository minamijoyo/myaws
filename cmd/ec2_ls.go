package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ec2LsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List EC2 instances",
	Run:   ec2.Ls,
}

func init() {
	ec2Cmd.AddCommand(ec2LsCmd)

	ec2LsCmd.Flags().BoolP("all", "a", false, "List all instances (by default, list running instances only)")
	ec2LsCmd.Flags().StringP("filter-tag", "t", "Name:",
		"Filter instances by tag, such as \"Name:app-production\". The value of tag is assumed to be a partial match",
	)
	ec2LsCmd.Flags().StringP("fields", "F", "InstanceId InstanceType PublicIpAddress PrivateIpAddress StateName LaunchTime Tag:Name", "Output fields list separated by space")

	viper.BindPFlag("ec2.ls.all", ec2LsCmd.Flags().Lookup("all"))
	viper.BindPFlag("ec2.ls.filter-tag", ec2LsCmd.Flags().Lookup("filter-tag"))
	viper.BindPFlag("ec2.ls.fields", ec2LsCmd.Flags().Lookup("fields"))
}
