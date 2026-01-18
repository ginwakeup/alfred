package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ginwakeup/alfred/cli/internal/docker"
	"gopkg.in/yaml.v3"
)

type SourceType string

type AlfredConfig struct {
	Project struct {
		Name    string `yaml:"name"`
		Compose string `yaml:"-"`
	} `yaml:"project"`

	Dependencies struct {
		RepositoryType string `yaml:"repository_type"`
		// This can either be a local file system path, or a GitHub repository, depending on RepositoryType.
		Location     string   `yaml:"location"`
		Dependencies []string `yaml:"dependencies"`
	}
	Network struct {
		Name string `yaml:"name"`
	} `yaml:"network"`

	// Internal - runtime
	Path     string `yaml:"-"`
	CacheDir string `yaml:"-"`
}

func (cfg *AlfredConfig) Init(path string) error {
	// LoadConfig alfred config data and unmarshal to AlfredConfig struct.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cfg.Path = path
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	// Create a cache directory
	configRootDir := filepath.Dir(cfg.Path)
	cacheDir := filepath.Join(configRootDir, ".alfred")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}
	cfg.CacheDir = cacheDir

	// Setup some initial values

	// Alfred assumes a docker-compose in the project-root, for the moment.
	cfg.Project.Compose = filepath.Join(filepath.Dir(path), "docker-compose.yaml")

	// Run some validation
	if err := cfg.Validate(); err != nil {
		return err
	}

	return nil
}

func (cfg *AlfredConfig) Create() error {
	// Create project dir if it does not exist.
	projectRoot := filepath.Dir(cfg.Path)
	_, err := os.Stat(projectRoot)
	if os.IsNotExist(err) {
		err = os.MkdirAll(projectRoot, 0755)
		if err != nil {
			return err
		}
	}

	// YAML data Marshall
	data, err := yaml.Marshal(&cfg)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(cfg.Path, data, 0644)
	if err != nil {
		panic(err)
	}
	return nil
}

func (cfg *AlfredConfig) Validate() error {
	if cfg.Network.Name == "" {
		return fmt.Errorf("No network name specified in Alfred Config.")
	}

	// Validate Project Compose for Alfred features
	if err := docker.Validate(cfg.Project.Compose); err != nil {
		return err
	}

	for _, system := range cfg.Dependencies.Dependencies {
		dockerComposePath := fmt.Sprintf("%s/%s/docker-compose.yaml", cfg.Dependencies.Location, system)
		if err := docker.Validate(dockerComposePath); err != nil {
			return err
		}
	}
	return nil
}

func (cfg *AlfredConfig) RunDependencies() error {
	for _, system := range cfg.Dependencies.Dependencies {
		depComposePath := fmt.Sprintf("%s/%s/docker-compose.yaml", cfg.Dependencies.Location, system)
		outTmpComposePath := filepath.Join(cfg.CacheDir, system, "docker-compose.yaml")

		// Before running, add a custom network and write tmp output yamls
		tmpComposePath := docker.GenerateTmpCompose(depComposePath, outTmpComposePath)
		fmt.Println("tmpComposePath:", tmpComposePath)
		if err := docker.Up(outTmpComposePath, cfg.Project.Name); err != nil {
			return err
		}
	}
	return nil
}

func LoadConfig(path string) (*AlfredConfig, error) {
	var cfg AlfredConfig

	if err := cfg.Init(path); err != nil {
		return nil, err
	}

	return &cfg, nil
}
