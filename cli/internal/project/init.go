package project

import (
	"path/filepath"

	"github.com/ginwakeup/alfred/cli/internal/config"
)

func InitAlfredProject(projectName string, dependenciesRoot string, projectRoot string, dependencies []string) {
	// Create Alfred Config YAML
	var cfg config.AlfredConfig

	cfg.Project.Name = projectName
	cfg.DependenciesRoot = dependenciesRoot
	cfg.Dependencies = dependencies
	cfg.Path = filepath.Join(projectRoot, "alfred.yaml")
	cfg.Network.Name = "alfred-dev"

	err := cfg.Create()
	if err != nil {
		return
	}
}
