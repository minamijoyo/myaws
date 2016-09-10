package myaws

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func UsageError(cmd *cobra.Command, msg interface{}) {
	fmt.Println("Error:", msg)
	cmd.Usage()
	os.Exit(-1)
}
