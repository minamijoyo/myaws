package cmd

import (
	"github.com/minamijoyo/myaws/myaws"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(newACMCmd())
}

func newACMCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acm",
		Short: "Manage ACM resources",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		newACMLsCmd(),
	)

	return cmd
}

func newACMLsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List certificates",
		RunE:  runACMLsCmd,
	}

	flags := cmd.Flags()
	flags.BoolP("verbose", "v", false, "Verbose")
	flags.BoolP("pending", "p", false, "List only pending certificates")

	viper.BindPFlag("acm.ls.verbose", flags.Lookup("verbose"))
	viper.BindPFlag("acm.ls.pending", flags.Lookup("pending"))

	return cmd
}

func runACMLsCmd(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return errors.Wrap(err, "newClient failed:")
	}

	options := myaws.ACMLsOptions{
		Verbose: viper.GetBool("acm.ls.verbose"),
		Pending: viper.GetBool("acm.ls.pending"),
	}

	return client.ACMLs(options)
}
