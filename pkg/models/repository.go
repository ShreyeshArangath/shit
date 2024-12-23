package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ShreyeshArangath/shit/pkg/utils"
	"gopkg.in/ini.v1"
)

// Repository represents a shit repository with its worktree and git directory.
type Repository struct {
	Worktree string    // Worktree is the directory where the working files are checked out.
	GitDir   string    // GitDir is the directory where the Git metadata is stored.
	Conf     *ini.File // Conf is the configuration for the repository.
}

// Create a new repository
func CreateRepository(path string, force bool) (*Repository, error) {
	// Check if the path is valid

	// Check if the path has a .git subdir (if not, throw an exception )
	gitfilepath := filepath.Join(path, utils.GIT_DIR_NAME)
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !(force || info.IsDir()) {
		return nil, &ShitException{Message: fmt.Sprintf("Not a git repository %s", path)}
	}
	// Read the .git config file
	cfg, err := ini.Load(utils.CONFIG_FILE_NAME)
	if err != nil {
		return nil, err
	}
	return &Repository{
		Worktree: path,
		GitDir:   gitfilepath,
		Conf:     cfg,
	}, nil
}

// Get the path under the .git directory
func (r *Repository) GetRepoPath(paths ...string) string {
	pathwithgitdir := append([]string{r.GitDir}, paths...)
	return filepath.Join(pathwithgitdir...)
}
