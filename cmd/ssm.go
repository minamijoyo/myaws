package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws"
)

func init() {
	RootCmd.AddCommand(newSSMCmd())
}

func newSSMCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssm",
		Short: "Manage SSM resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newSSMParameterCmd(),
	)

	return cmd
}

func newSSMParameterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parameter",
		Short: "Manage SSM parameter resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newSSMParameterPutCmd(),
		newSSMParameterGetCmd(),
		newSSMParameterLsCmd(),
		newSSMParameterEnvCmd(),
		newSSMParameterDelCmd(),
	)

	return cmd
}

func newSSMParameterPutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put NAME VALUE",
		Short: "Put SSM parameter",
		RunE:  runSSMParameterPutCmd,
	}

	flags := cmd.Flags()
	flags.StringP("key-id", "k", "", "KMS key ID or alias")

	viper.BindPFlag("ssm.parameter.put.key-id", flags.Lookup("key-id"))

	return cmd
}

func runSSMParameterPutCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) != 2 {
		return errors.New("NAME and VALUE are required")
	}

	options := myaws.SSMParameterPutOptions{
		Name:  args[0],
		Value: args[1],
		KeyID: viper.GetString("ssm.parameter.put.key-id"),
	}

	return client.SSMParameterPut(options)
}

func newSSMParameterGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get NAME [...]",
		Short: "Get SSM parameter",
		RunE:  runSSMParameterGetCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("with-decryption", "d", true, "with KMS decryption")

	viper.BindPFlag("ssm.parameter.get.with-decryption", flags.Lookup("with-decryption"))
	return cmd
}

func runSSMParameterGetCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("NAME is required")
	}

	names := aws.StringSlice(args)
	options := myaws.SSMParameterGetOptions{
		Names:          names,
		WithDecryption: viper.GetBool("ssm.parameter.get.with-decryption"),
	}

	return client.SSMParameterGet(options)
}

func newSSMParameterLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List SSM parameters",
		RunE:  runSSMParameterLsCmd,
	}

	flags := cmd.Flags()
	flags.StringP("name", "n", "",
		"Filter parameters by Name, such as foo.dev. The value of tag is assumed to be a prefix match",
	)

	viper.BindPFlag("ssm.parameter.ls.name", flags.Lookup("name"))
	return cmd
}

func runSSMParameterLsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.SSMParameterLsOptions{
		Name: viper.GetString("ssm.parameter.ls.name"),
	}

	return client.SSMParameterLs(options)
}

func newSSMParameterEnvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "env NAME",
		Short: "Print SSM parameters as a list of environment variables",
		RunE:  runSSMParameterEnvCmd,
	}

	return cmd
}

func runSSMParameterEnvCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("NAME is required")
	}

	options := myaws.SSMParameterEnvOptions{
		Name: args[0],
	}

	return client.SSMParameterEnv(options)
}

func newSSMParameterDelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del NAME",
		Short: "Delete SSM parameter",
		RunE:  runSSMParameterDelCmd,
	}

	return cmd
}

func runSSMParameterDelCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	if len(args) == 0 {
		return errors.New("NAME is required")
	}

	options := myaws.SSMParameterDelOptions{
		Name: args[0],
	}

	return client.SSMParameterDel(options)
}
