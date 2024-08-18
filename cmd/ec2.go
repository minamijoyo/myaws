package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws"
)

func init() {
	RootCmd.AddCommand(newEC2Cmd())
}

func newEC2Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ec2",
		Short: "Manage EC2 resources",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Help() // nolint: errcheck
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
	flags.StringP("filter-tag", "t", "",
		"Filter instances by tag, such as \"Name:app-production\". The value of tag is assumed to be a partial match",
	)
	flags.StringP("fields", "F", "InstanceId InstanceType PublicIpAddress PrivateIpAddress AvailabilityZone StateName LaunchTime Tag:Name", "Output fields list separated by space")

	viper.BindPFlag("ec2.ls.all", flags.Lookup("all"))               // nolint: errcheck
	viper.BindPFlag("ec2.ls.quiet", flags.Lookup("quiet"))           // nolint: errcheck
	viper.BindPFlag("ec2.ls.filter-tag", flags.Lookup("filter-tag")) // nolint: errcheck
	viper.BindPFlag("ec2.ls.fields", flags.Lookup("fields"))         // nolint: errcheck

	return cmd
}

func runEC2LsCmd(_ *cobra.Command, _ []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.EC2LsOptions{
		All:       viper.GetBool("ec2.ls.all"),
		Quiet:     viper.GetBool("ec2.ls.quiet"),
		FilterTag: viper.GetString("ec2.ls.filter-tag"),
		Fields:    viper.GetStringSlice("ec2.ls.fields"),
	}

	return client.EC2Ls(options)
}

func newEC2StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start INSTANCE_ID [...]",
		Short: "Start EC2 instances",
		RunE:  runEC2StartCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("wait", "w", false, "Wait until instance running")

	viper.BindPFlag("ec2.start.wait", flags.Lookup("wait")) // nolint: errcheck

	return cmd
}

func runEC2StartCmd(_ *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("INSTANCE_ID is required")
	}
	instanceIds := aws.StringSlice(args)

	options := myaws.EC2StartOptions{
		InstanceIds: instanceIds,
		Wait:        viper.GetBool("ec2.start.wait"),
	}

	return client.EC2Start(options)
}

func newEC2StopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop INSTANCE_ID [...]",
		Short: "Stop EC2 instances",
		RunE:  runEC2StopCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("wait", "w", false, "Wait until instance stopped")

	viper.BindPFlag("ec2.stop.wait", flags.Lookup("wait")) // nolint: errcheck

	return cmd
}

func runEC2StopCmd(_ *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("INSTANCE_ID is required")
	}
	instanceIds := aws.StringSlice(args)

	options := myaws.EC2StopOptions{
		InstanceIds: instanceIds,
		Wait:        viper.GetBool("ec2.stop.wait"),
	}

	return client.EC2Stop(options)
}

func newEC2SSHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh [USER@]INSTANCE_NAME [COMMAND...]",
		Short: "SSH to EC2 instances",
		RunE:  runEC2SSHCmd,
	}

	flags := cmd.Flags()
	flags.StringP("login-name", "l", "", "Login username")
	flags.StringP("identity-file", "i", "~/.ssh/id_rsa", "SSH private key file")
	flags.BoolP("private", "", false, "Use private IP to connect")

	viper.BindPFlag("ec2.ssh.login-name", flags.Lookup("login-name"))       // nolint: errcheck
	viper.BindPFlag("ec2.ssh.identity-file", flags.Lookup("identity-file")) // nolint: errcheck
	viper.BindPFlag("ec2.ssh.private", flags.Lookup("private"))             // nolint: errcheck

	return cmd
}

func runEC2SSHCmd(_ *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("Instance name is required")
	}

	var loginName, instanceName string
	if strings.Contains(args[0], "@") {
		// parse loginName@instanceName format
		splitted := strings.SplitN(args[0], "@", 2)
		loginName, instanceName = splitted[0], splitted[1]
	} else {
		loginName = viper.GetString("ec2.ssh.login-name")
		instanceName = args[0]
	}

	filterTag := "Name:" + instanceName

	var command string
	if len(args) >= 2 {
		command = strings.Join(args[1:], " ")
	}
	options := myaws.EC2SSHOptions{
		FilterTag:    filterTag,
		LoginName:    loginName,
		IdentityFile: viper.GetString("ec2.ssh.identity-file"),
		Private:      viper.GetBool("ec2.ssh.private"),
		Command:      command,
	}

	return client.EC2SSH(options)
}
