package core

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GetRepo clones or pulls a repo into cacheDir.
// Uses a token from the environment variable GITHUB_TOKEN.
func GetRepo(repoURL string, branch string, cacheDir string) (string, error) {
	hash := sha1.Sum([]byte(repoURL)) // returns [20]byte

	// Convert to hex string
	hashStr := fmt.Sprintf("%x", hash)
	repoPath := filepath.Join(cacheDir, hashStr)

	// Get token from env
	username := os.Getenv("GITHUB_USERNAME")
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return "", fmt.Errorf("GITHUB_TOKEN not set")
	}

	auth := &http.BasicAuth{
		Username: username, // can be anything except empty string
		Password: token,    // your personal access token
	}

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		// Clone if repo doesn't exist
		_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:           repoURL,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
			SingleBranch:  true,
			Auth:          auth,
		})
		if err != nil {
			return "", err
		}
	} else {
		// Pull latest changes
		r, err := git.PlainOpen(repoPath)
		if err != nil {
			return "", err
		}

		w, err := r.Worktree()
		if err != nil {
			return "", err
		}

		err = w.Pull(&git.PullOptions{
			RemoteName:    "origin",
			ReferenceName: plumbing.NewBranchReferenceName(branch),
			SingleBranch:  true,
			Auth:          auth,
		})

		// git.NoErrAlreadyUpToDate is not fatal
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return "", err
		}
	}

	return repoPath, nil
}
