package myaws

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// UsageError prints a error message and show usage, then exit with return -1.
func UsageError(cmd *cobra.Command, msg interface{}) {
	fmt.Println("Error:", msg)
	cmd.Usage()
	os.Exit(-1)
}
