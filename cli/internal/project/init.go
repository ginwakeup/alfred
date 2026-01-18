package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ginwakeup/alfred/cli/internal/config"
)

func InitAlfredProject(projectName string, projectRoot string, dependenciesRoot string, dependenciesRepoType string, dependencies []string) error {
	// Create Alfred Config YAML
	var cfg config.AlfredConfig

	cfg.Project.Name = projectName
	cfg.Dependencies.Location = dependenciesRoot
	cfg.Dependencies.Dependencies = dependencies
	cfg.Dependencies.RepositoryType = dependenciesRepoType
	cfg.Path = filepath.Join(projectRoot, "alfred.yaml")
	cfg.Network.Name = "alfred-dev"

	_, err := os.Stat(cfg.Path)
	if os.IsNotExist(err) {
		err := cfg.Create()
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Alfred Project already exists, it won't be created.")
	}

	if err != nil {
		// some other error, e.g., permission denied
		fmt.Println("Error checking directory:", err)
		return err
	}
	return nil
}
