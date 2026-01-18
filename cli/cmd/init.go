package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ginwakeup/alfred/cli/internal/project"
	"github.com/spf13/cobra"
)

type InitConfig struct {
	ProjectRoot      string
	ProjectName      string
	ProjectCompose   string
	DependenciesRoot string
	Dependencies     []string
}

var (
	projectName      string
	projectRoot      string
	dependenciesRoot string
	dependenciesRaw  string
)

func resolveInitConfig() (*InitConfig, error) {
	if projectName == "" {
		return nil, fmt.Errorf("--project-name is required")
	}

	if projectRoot == "" {
		return nil, fmt.Errorf("--project-root is required")
	}

	depsRoot := dependenciesRoot
	if depsRoot == "" {
		depsRoot = os.Getenv("DEPENDENCIES_ROOT")
	}

	if depsRoot == "" {
		return nil, fmt.Errorf(
			"--dependencies-root is required unless DEPENDENCIES_ROOT is set",
		)
	}

	var deps []string
	if dependenciesRaw != "" {
		deps = strings.Split(dependenciesRaw, ",")
		for i := range deps {
			deps[i] = strings.TrimSpace(deps[i])
		}
	}

	return &InitConfig{
		ProjectName:      projectName,
		ProjectRoot:      projectRoot,
		DependenciesRoot: depsRoot,
		Dependencies:     deps,
	}, nil
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := resolveInitConfig()
		if err != nil {
			return err
		}

		return runInit(cfg)
	},
}

func runInit(cfg *InitConfig) error {
	fmt.Println("Initializing project")
	fmt.Println("Project Name:", cfg.ProjectName)
	fmt.Println("Project root:", cfg.ProjectRoot)
	fmt.Println("Dependencies root:", cfg.DependenciesRoot)
	fmt.Println("Dependencies:", cfg.Dependencies)

	project.InitAlfredProject(cfg.ProjectName, cfg.DependenciesRoot, cfg.ProjectRoot, cfg.Dependencies)
	return nil
}

func Init() *cobra.Command {
	initCmd.Flags().StringVar(
		&projectName,
		"project-name",
		"",
		"Name of the project (required)",
	)

	initCmd.Flags().StringVar(
		&projectRoot,
		"project-root",
		"",
		"Path to the project root (required)",
	)

	initCmd.Flags().StringVar(
		&dependenciesRoot,
		"dependencies-root",
		"",
		"Path to the dependencies root (or set DEPENDENCIES_ROOT)",
	)

	initCmd.Flags().StringVar(
		&dependenciesRaw,
		"dependencies",
		"",
		"Comma-separated list of dependencies (optional)",
	)

	_ = initCmd.MarkFlagRequired("project-root")

	return initCmd
}
