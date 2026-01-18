package main

import (
	"fmt"
	"os"

	"github.com/ginwakeup/alfred/cli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "",
		Short: "Bootstrap a dev environment for your application",
	}

	root.AddCommand(cmd.Run())
	root.AddCommand(cmd.Init())

	if err := root.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
