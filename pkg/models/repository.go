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
	gitfilepath := filepath.Join(path, utils.GIT_DIR_NAME)
	isdir, _ := utils.IsDir(gitfilepath)
	if !(force || isdir) {
		return nil, &ShitException{Message: fmt.Sprintf("Not a git repository %s", path)}
	}
	repository := &Repository{
		Worktree: path,
		GitDir:   gitfilepath,
	}

	pathtoconfig, _ := repository.RepoFile(false, utils.CONFIG_FILE_NAME)
	exists, _ := utils.PathExists(pathtoconfig)
	if exists && pathtoconfig != "" {
		// Load the configuration file
		cfg, err := ini.Load(pathtoconfig)
		if err != nil {
			return nil, err
		}
		repository.Conf = cfg
	} else if !force {
		return nil, &ShitException{Message: fmt.Sprintf("Config file does not exist %s", pathtoconfig)}
	}

	if !force {
		version := repository.Conf.Section("core").Key("repositoryformatversion").String()
		if version != "0" {
			return nil, &ShitException{Message: fmt.Sprintf("Unsupported repositoryformatversion %s", version)}
		}
	}
	return repository, nil
}

// Get the path under the .git directory
func (r *Repository) GetRepoPath(path ...string) string {
	pathwithgitdir := append([]string{r.GitDir}, path...)
	return filepath.Join(pathwithgitdir...)
}

// RepoDir returns the directory path for the repository, creating it if necessary.
// If the directory already exists, it verifies that the path is indeed a directory.
// If the directory does not exist and mkdir is true, it creates the directory with
// the specified path. If mkdir is false, it returns an error indicating that the
// path does not exist.
//
// Parameters:
//
//	mkdir - a boolean indicating whether to create the directory if it does not exist
//	path  - variadic string arguments representing the path components to the directory
//
// Returns:
//
//	string - the path to the repository directory
//	error  - an error if the path is not a directory, does not exist, or if there
//	         was an issue creating the directory
func (r *Repository) RepoDir(mkdir bool, path ...string) (string, error) {
	pathtodir := r.GetRepoPath(path...)
	exists, err := utils.PathExists(pathtodir)
	if err != nil {
		return "", err
	}
	if exists {
		isdir, err := utils.IsDir(pathtodir)
		if err != nil {
			return "", err
		}
		if isdir {
			return pathtodir, nil
		}
		return "", &ShitException{Message: fmt.Sprintf("Not a directory %s", pathtodir)}
	}

	if mkdir {
		if err := os.MkdirAll(pathtodir, 0755); err != nil {
			return "", err
		}
		return pathtodir, nil
	}
	return "", &ShitException{Message: fmt.Sprintf("Path does not exist %s", pathtodir)}
}

// RepoFile creates the repository directory if `mkdir` is true and returns the
// full path to the repository file specified by the `path` arguments.
// It returns an error if the directory creation or path retrieval fails.
//
// Parameters:
// - mkdir: A boolean indicating whether to create the directory if it does not exist.
// - path: A variadic string slice representing the path components of the repository file.
//
// Returns:
// - A string representing the full path to the repository file.
// - An error if the directory creation or path retrieval fails.
func (r *Repository) RepoFile(mkdir bool, path ...string) (string, error) {
	_, err := r.RepoDir(mkdir, path[:len(path)-1]...)
	if err != nil {
		return "", err
	}
	return r.GetRepoPath(path...), nil
}

// RepoFind recursively finds the ".git" directory starting from the given path.
func RepoFind(path string, required bool) (*Repository, error) {
	// Convert the path to an absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// Check if ".git" exists in the current directory
	gitPath := filepath.Join(absPath, ".git")
	isDir, err := utils.IsDir(gitPath)
	if isDir {
		return CreateRepository(absPath, false)
	}

	// Determine the parent directory
	parent := filepath.Dir(absPath)
	if parent == absPath {
		// If we've reached the root directory
		if required {
			return nil, &ShitException{Message: fmt.Sprintf("no git directory found in %s", path)}
		}
		return nil, nil
	}
	return RepoFind(parent, required)
}
