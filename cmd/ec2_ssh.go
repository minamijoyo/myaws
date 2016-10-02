package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ec2SshCmd = &cobra.Command{
	Use:   "ssh INSTANCE_ID",
	Short: "SSH to EC2 instance",
	Run:   ec2.Ssh,
}

func init() {
	ec2Cmd.AddCommand(ec2SshCmd)
	ec2SshCmd.Flags().StringP("login-name", "l", "", "Login username")
	ec2SshCmd.Flags().StringP("identity-file", "i", "~/.ssh/id_rsa", "SSH private key file")

	viper.BindPFlag("ec2.ssh.login-name", ec2SshCmd.Flags().Lookup("login-name"))
	viper.BindPFlag("ec2.ssh.identity-file", ec2SshCmd.Flags().Lookup("identity-file"))
}
