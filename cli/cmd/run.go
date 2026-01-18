package cmd

import (
	"path/filepath"

	"github.com/ginwakeup/alfred/cli/internal/config"
	"github.com/ginwakeup/alfred/cli/internal/docker"
	"github.com/spf13/cobra"
)

func Run() *cobra.Command {
	return &cobra.Command{
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
				return docker.Up(projectTmpComposeOut, cfg.Project.Name)
			}
			return err
		},
	}
}
