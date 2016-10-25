package cmd

import (
	"github.com/minamijoyo/myaws/myaws/ecr"
	"github.com/spf13/cobra"
)

var ecrGetLoginCmd = &cobra.Command{
	Use:   "get-login",
	Short: "Get docker login command for ECR",
	Run:   ecr.GetLogin,
}

func init() {
	ecrCmd.AddCommand(ecrGetLoginCmd)
}
