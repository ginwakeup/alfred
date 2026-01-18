package main

import (
	"fmt"
	"os"

	"github.com/ginwakeup/alfred/cli/cmd"
	"github.com/ginwakeup/alfred/cli/internal/core"
	"github.com/ginwakeup/alfred/cli/internal/core/types"
	"github.com/spf13/cobra"
)

func main() {
	// First create or retrieve the Cache Directory.
	dir, err := core.CreateCacheDir(".alfred")
	if err != nil {
		return
	}

	alfredRunTimeCfg := types.AlfredRunTimeConfig{
		CacheDir: dir,
	}

	root := &cobra.Command{
		Use:   "",
		Short: "Bootstrap a dev environment for your application",
	}

	root.AddCommand(cmd.Run(&alfredRunTimeCfg))
	root.AddCommand(cmd.Init())

	if err := root.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
