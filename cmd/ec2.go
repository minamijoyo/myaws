package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws/ec2"
)

func init() {
	RootCmd.AddCommand(newEC2Cmd())
}

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
		RunE:  runEC2LsCmd,
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

func runEC2LsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := ec2.LsOptions{
		All:       viper.GetBool("ec2.ls.all"),
		Quiet:     viper.GetBool("ec2.ls.quiet"),
		FilterTag: viper.GetString("ec2.ls.filter-tag"),
		Fields:    viper.GetStringSlice("ec2.ls.fields"),
	}

	return ec2.Ls(client, options)
}

func newEC2StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start INSTANCE_ID [...]",
		Short: "Start EC2 instances",
		RunE:  runEC2StartCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("wait", "w", false, "Wait until instance running")

	viper.BindPFlag("ec2.start.wait", flags.Lookup("wait"))

	return cmd
}

func runEC2StartCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("INSTANCE_ID is required")
	}
	instanceIds := aws.StringSlice(args)

	options := ec2.StartOptions{
		InstanceIds: instanceIds,
		Wait:        viper.GetBool("ec2.start.wait"),
	}

	return ec2.Start(client, options)
}

func newEC2StopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop INSTANCE_ID [...]",
		Short: "Stop EC2 instances",
		RunE:  runEC2StopCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("wait", "w", false, "Wait until instance stopped")

	viper.BindPFlag("ec2.stop.wait", flags.Lookup("wait"))

	return cmd
}

func runEC2StopCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("INSTANCE_ID is required")
	}
	instanceIds := aws.StringSlice(args)

	options := ec2.StopOptions{
		InstanceIds: instanceIds,
		Wait:        viper.GetBool("ec2.stop.wait"),
	}

	return ec2.Stop(client, options)
}

func newEC2SSHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh INSTANCE_ID",
		Short: "SSH to EC2 instance",
		RunE:  ec2.SSH,
	}

	flags := cmd.Flags()
	flags.StringP("login-name", "l", "", "Login username")
	flags.StringP("identity-file", "i", "~/.ssh/id_rsa", "SSH private key file")

	viper.BindPFlag("ec2.ssh.login-name", flags.Lookup("login-name"))
	viper.BindPFlag("ec2.ssh.identity-file", flags.Lookup("identity-file"))

	return cmd
}
