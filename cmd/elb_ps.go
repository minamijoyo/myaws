package cmd

import (
	"github.com/minamijoyo/myaws/myaws/elb"
	"github.com/spf13/cobra"
)

var elbPsCmd = &cobra.Command{
	Use:   "ps ELB_NAME",
	Short: "Show ELB instances",
	Run:   elb.Ps,
}

func init() {
	elbCmd.AddCommand(elbPsCmd)
}
