package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ginwakeup/alfred/cli/internal/project"
	"github.com/spf13/cobra"
)

type InitConfig struct {
	ProjectRoot          string   // root location of the Alfred project to create.
	ProjectName          string   // name of the Alfred project to create
	DependenciesRepoType string   // git or local
	DependenciesLocation string   // a git repository url or a local file system path
	Dependencies         []string // list of dependencies to resolve
}

var (
	projectName          string
	projectRoot          string
	dependenciesRepoType string
	dependenciesLocation string
	dependenciesRaw      string
)

func resolveInitConfig() (*InitConfig, error) {
	if projectName == "" {
		return nil, fmt.Errorf("--project-name is required")
	}

	if projectRoot == "" {
		return nil, fmt.Errorf("--project-root is required")
	}

	if dependenciesRepoType == "" {
		return nil, fmt.Errorf("--dependencies-repo-type is required")
	}

	if dependenciesLocation == "" {
		return nil, fmt.Errorf("--dependencies-location is required")
	}

	if dependenciesRaw == "" {
		return nil, fmt.Errorf("--dependencies are required")
	}

	depsLocation := dependenciesLocation
	if depsLocation == "" {
		depsLocation = os.Getenv("DEPENDENCIES_LOCATION")
	}

	if depsLocation == "" {
		return nil, fmt.Errorf(
			"--dependencies-location is required unless DEPENDENCIES_LOCATION is set",
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
		ProjectName:          projectName,
		ProjectRoot:          projectRoot,
		DependenciesLocation: depsLocation,
		DependenciesRepoType: dependenciesRepoType,
		Dependencies:         deps,
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
	fmt.Println("Dependencies Repository Type:", cfg.DependenciesRepoType)
	fmt.Println("Dependencies Location:", cfg.DependenciesLocation)
	fmt.Println("Dependencies:", cfg.Dependencies)

	project.InitAlfredProject(cfg.ProjectName, cfg.ProjectRoot, cfg.DependenciesLocation, cfg.DependenciesRepoType, cfg.Dependencies)
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
		&dependenciesRepoType,
		"dependencies-repo-type",
		"",
		"Type of repository for dependencies: git or local",
	)

	initCmd.Flags().StringVar(
		&dependenciesLocation,
		"dependencies-location",
		"",
		"Path to the dependencies location (or set DEPENDENCIES_LOCATION)",
	)

	initCmd.Flags().StringVar(
		&dependenciesRaw,
		"dependencies",
		"",
		"Comma-separated list of dependencies",
	)

	_ = initCmd.MarkFlagRequired("project-name")
	_ = initCmd.MarkFlagRequired("project-name")
	_ = initCmd.MarkFlagRequired("dependencies-repo-type")
	_ = initCmd.MarkFlagRequired("dependencies-location")
	_ = initCmd.MarkFlagRequired("dependencies")

	return initCmd
}
