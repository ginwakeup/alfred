package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ginwakeup/alfred/cli/internal/config"
	"github.com/ginwakeup/alfred/cli/internal/docker"
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "dev",
		Short: "Bootstrap a dev environment for your application",
	}

	run := &cobra.Command{
		Use:   "run <alfred.yaml>",
		Args:  cobra.ExactArgs(1),
		Short: "Run development environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			// LoadConfig Alfred Configuration.
			cfg, err := config.LoadConfig(args[0])
			if err != nil {
				return err
			}

			// Create alfred-dev network
			if err := docker.CreateDevNetwork(cfg.Network.Name); err != nil {
				return err
			}

			// Run all dependencies.
			if err := cfg.RunDependencies(); err != nil {
				return err
			}

			// Apply overrides to Project compose and start
			projectTmpComposeOut := filepath.Join(cfg.CacheDir, "docker-compose.yaml")
			err = docker.GenerateTmpCompose(cfg.Project.Compose, projectTmpComposeOut)
			if err == nil {
				return docker.Up(projectTmpComposeOut)
			}
			return err
		},
	}

	root.AddCommand(run)

	if err := root.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
