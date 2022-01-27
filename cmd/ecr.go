package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/myaws"
)

func init() {
	RootCmd.AddCommand(newECRCmd())
}

func newECRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecr",
		Short: "Manage ECR resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help() // nolint: errcheck
		},
	}

	cmd.AddCommand(
		newECRGetLoginCmd(),
	)

	return cmd
}

func newECRGetLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-login",
		Short: "Get docker login command for ECR",
		RunE:  runECRGetLoginCmd,
	}

	flags := cmd.Flags()
	flags.StringSliceP("registry-ids", "r", []string{}, "A list of AWS account IDs")

	viper.BindPFlag("ecr.get-login.registry-ids", flags.Lookup("registry-ids")) // nolint: errcheck

	return cmd
}

func runECRGetLoginCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	registryIds := aws.StringSlice(viper.GetStringSlice("ecr.get-login.registry-ids"))
	options := myaws.ECRGetLoginOptions{
		RegistryIds: registryIds,
	}

	return client.ECRGetLogin(options)
}
