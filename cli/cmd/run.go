package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/ginwakeup/alfred/cli/internal/config"
	"github.com/ginwakeup/alfred/cli/internal/core/types"
	"github.com/ginwakeup/alfred/cli/internal/docker"
	"github.com/spf13/cobra"
)

func Run(alfredRunTimeCfg *types.AlfredRunTimeConfig) *cobra.Command {
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
			if err := cfg.RunDependencies(alfredRunTimeCfg); err != nil {
				return err
			}

			// Apply overrides to Project compose and start, only if a compose for the project exists.
			projectTmpComposeOut := filepath.Join(cfg.CacheDir, "docker-compose.yaml")

			if _, err := os.Stat(projectTmpComposeOut); errors.Is(err, os.ErrNotExist) {
				return nil
			}

			if cfg.Project.Compose != "" {
				err = docker.GenerateOverriddenCompose(cfg.Project.Compose, projectTmpComposeOut)
				if err == nil {
					return docker.Up(projectTmpComposeOut, cfg.Project.Name)
				}
				return err
			}
			return nil
		},
	}
}
