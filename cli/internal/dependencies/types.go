package dependencies

import (
	"fmt"

	"github.com/ginwakeup/alfred/cli/internal/core"
)
import "github.com/ginwakeup/alfred/cli/internal/core/types"

type Dependencies struct {
	RepositoryType string `yaml:"repository_type"`
	// This can either be a local file system path, or a GitHub repository, depending on RepositoryType.
	Location     string   `yaml:"location"`
	Dependencies []string `yaml:"dependencies"`
}

// Resolve dependencies location in case user selects "git", which means: pull repo locally and return all compose paths.
func (deps *Dependencies) ResolveDependenciesLocation(alfredRunTimeCfg *types.AlfredRunTimeConfig) ([]string, error) {
	var composePaths []string
	var dependenciesRootPath string

	switch deps.RepositoryType {
	case "git":
		// Pull the repo in the cache location.
		repoPath, err := core.GetRepo(deps.Location, "main", alfredRunTimeCfg.CacheDir)
		if err != nil {
			return nil, err
		}

		dependenciesRootPath = repoPath
	case "local":
		dependenciesRootPath = deps.Location
	}

	// Loop all dependencies, compute their compose paths and return them.
	for _, system := range deps.Dependencies {
		depComposePath := fmt.Sprintf("%s/%s/docker-compose.yaml", dependenciesRootPath, system)
		composePaths = append(composePaths, depComposePath)
	}

	return composePaths, nil
}
