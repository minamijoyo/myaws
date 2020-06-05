package cmd

import (
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
			cmd.Help()
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
	flags.StringP("registry-ids", "r", "", "Amazon ECR registries ID")

	viper.BindPFlag("ecr.parameter.get-login.registry-ids", flags.Lookup("registry-ids"))

	return cmd
}

func runECRGetLoginCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.ECRGetLoginOptions{
		RegistryIds: viper.GetString("ecr.parameter.get-login.registry-ids"),
	}

	return client.ECRGetLogin(options)
}
