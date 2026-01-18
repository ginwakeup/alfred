package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ginwakeup/alfred/cli/internal/docker"
	"gopkg.in/yaml.v3"
)

type DevConfig struct {
	Project struct {
		Name    string `yaml:"name"`
		Compose string `yaml:"compose"`
	} `yaml:"project"`

	Dependencies     []string `yaml:"dependencies"`
	DependenciesRoot string   `yaml:"dependencies_root"`
	Network          struct {
		Name string `yaml:"name"`
	} `yaml:"network"`

	ConfigPath string
	CacheDir   string
}

func (cfg *DevConfig) Init(path string) error {
	// LoadConfig alfred config data and unmarshal to DevConfig struct.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	cfg.ConfigPath = path
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	// Create a cache directory
	configRootDir := filepath.Dir(cfg.ConfigPath)
	cacheDir := filepath.Join(configRootDir, ".alfred")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}
	cfg.CacheDir = cacheDir

	// Setup some initial values

	// If config.app.compose is absolute, leave it as is.
	// If not, resolve it based on alfred.yaml location.
	if !filepath.IsAbs(cfg.Project.Compose) {
		cfg.Project.Compose = filepath.Join(filepath.Dir(path), cfg.Project.Compose)
	}

	// Run some validation
	if err := cfg.Validate(); err != nil {
		return err
	}

	return nil
}

func (cfg *DevConfig) Validate() error {
	if cfg.Network.Name == "" {
		return fmt.Errorf("No network name specified in Alfred Config.")
	}

	// Validate Project Compose for Alfred features
	if err := docker.Validate(cfg.Project.Compose); err != nil {
		return err
	}

	for _, system := range cfg.Dependencies {
		dockerComposePath := fmt.Sprintf("%s/%s/docker-compose.yaml", cfg.DependenciesRoot, system)
		if err := docker.Validate(dockerComposePath); err != nil {
			return err
		}
	}
	return nil
}

func (cfg *DevConfig) RunDependencies() error {
	for _, system := range cfg.Dependencies {
		depComposePath := fmt.Sprintf("%s/%s/docker-compose.yaml", cfg.DependenciesRoot, system)
		outTmpComposePath := filepath.Join(cfg.CacheDir, system, "docker-compose.yaml")

		// Before running, add a custom network and write tmp output yamls
		tmpComposePath := docker.GenerateTmpCompose(depComposePath, outTmpComposePath)
		fmt.Println("tmpComposePath:", tmpComposePath)
		if err := docker.Up(outTmpComposePath); err != nil {
			return err
		}
	}
	return nil
}

func LoadConfig(path string) (*DevConfig, error) {
	var cfg DevConfig

	if err := cfg.Init(path); err != nil {
		return nil, err
	}

	return &cfg, nil
}
