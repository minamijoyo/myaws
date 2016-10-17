package cmd

import (
	"github.com/minamijoyo/myaws/myaws/elb"
	"github.com/spf13/cobra"
)

var elbLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List ELB instances",
	Run:   elb.Ls,
}

func init() {
	elbCmd.AddCommand(elbLsCmd)
}
