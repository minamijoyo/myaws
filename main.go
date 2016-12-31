package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/minamijoyo/myaws/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		if viper.GetBool("debug") {
			fmt.Printf("%+v\n", err)
		} else {
			fmt.Printf("%v\n", err)
		}
		os.Exit(1)
	}
}
