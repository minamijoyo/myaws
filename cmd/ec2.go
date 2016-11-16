package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws/ec2"
)

func newEC2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ec2",
		Short: "Manage EC2 resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newEC2LsCmd(),
		newEC2StartCmd(),
		newEC2StopCmd(),
		newEC2SSHCmd(),
	)

	return cmd
}

func newEC2LsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List EC2 instances",
		RunE:  ec2.Ls,
	}

	flags := cmd.Flags()
	flags.BoolP("all", "a", false, "List all instances (by default, list running instances only)")
	flags.BoolP("quiet", "q", false, "Only display InstanceIDs")
	flags.StringP("filter-tag", "t", "Name:",
		"Filter instances by tag, such as \"Name:app-production\". The value of tag is assumed to be a partial match",
	)
	flags.StringP("fields", "F", "InstanceId InstanceType PublicIpAddress PrivateIpAddress StateName LaunchTime Tag:Name", "Output fields list separated by space")

	viper.BindPFlag("ec2.ls.all", flags.Lookup("all"))
	viper.BindPFlag("ec2.ls.quiet", flags.Lookup("quiet"))
	viper.BindPFlag("ec2.ls.filter-tag", flags.Lookup("filter-tag"))
	viper.BindPFlag("ec2.ls.fields", flags.Lookup("fields"))

	return cmd
}

func newEC2StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start INSTANCE_ID [...]",
		Short: "Start EC2 instances",
		RunE:  ec2.Start,
	}

	flags := cmd.Flags()
	flags.BoolP("wait", "w", false, "Wait until instance running")

	viper.BindPFlag("ec2.start.wait", flags.Lookup("wait"))

	return cmd
}

func newEC2StopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop INSTANCE_ID [...]",
		Short: "Stop EC2 instances",
		RunE:  ec2.Stop,
	}

	flags := cmd.Flags()
	flags.BoolP("wait", "w", false, "Wait until instance stopped")

	viper.BindPFlag("ec2.stop.wait", flags.Lookup("wait"))

	return cmd
}

func newEC2SSHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh INSTANCE_ID",
		Short: "SSH to EC2 instance",
		Run:   ec2.SSH,
	}

	flags := cmd.Flags()
	flags.StringP("login-name", "l", "", "Login username")
	flags.StringP("identity-file", "i", "~/.ssh/id_rsa", "SSH private key file")

	viper.BindPFlag("ec2.ssh.login-name", flags.Lookup("login-name"))
	viper.BindPFlag("ec2.ssh.identity-file", flags.Lookup("identity-file"))

	return cmd
}
