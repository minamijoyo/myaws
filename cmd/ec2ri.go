package cmd

import (
	"github.com/minamijoyo/myaws/myaws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(newEC2RICmd())
}

func newEC2RICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ec2ri",
		Short: "Manage EC2 Reserved Instance resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newEC2RILsCmd(),
	)

	return cmd
}

func newEC2RILsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List EC2 Reserved Instances",
		RunE:  runEC2RILsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("all", "a", false, "List all reserved instances (by default, list active reserved instances only)")
	flags.BoolP("quiet", "q", false, "Only display ReservedInstanceIDs")
	flags.StringP("fields", "F", "ReservedInstancesId State Scope AvailabilityZone InstanceType InstanceCount Duration Start End", "Output fields list separated by space")

	viper.BindPFlag("ec2ri.ls.all", flags.Lookup("all"))
	viper.BindPFlag("ec2ri.ls.quiet", flags.Lookup("quiet"))
	viper.BindPFlag("ec2ri.ls.fields", flags.Lookup("fields"))

	return cmd
}

func runEC2RILsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.EC2RILsOptions{
		All:    viper.GetBool("ec2ri.ls.all"),
		Quiet:  viper.GetBool("ec2ri.ls.quiet"),
		Fields: viper.GetStringSlice("ec2ri.ls.fields"),
	}

	return client.EC2RILs(options)
}
