package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ShreyeshArangath/shit/pkg/models"
	"github.com/ShreyeshArangath/shit/pkg/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var log = logrus.New()

// .git/objects/ : the object store, which we’ll introduce in the next section.
// .git/refs/ the reference store, which we’ll discuss a bit later. It contains two subdirectories, heads and tags.
// .git/HEAD, a reference to the current HEAD (more on that later!)
// .git/config, the repository’s configuration file.
// .git/description, holds a free-form description of this repository’s contents, for humans, and is rarely used.
func Init(path string) (bool, error) {
	// Create the new shit repository
	repo, err := models.CreateRepository(path, true)
	if err != nil {
		return false, err
	}

	// Ensure the worktree path exists and is a directory
	workTreeExists, err := utils.IsDir(repo.Worktree)
	if err != nil {
		return false, err
	}

	if workTreeExists {
		workTreeIsDir, err := utils.IsDir(repo.Worktree)
		if err != nil {
			return false, err
		}
		if !workTreeIsDir {
			err = &models.ShitException{Message: fmt.Sprintf("Path %s is not a directory", repo.Worktree)}
			return false, err
		}

		gitDirExists, err := utils.PathExists(repo.GitDir)
		if err != nil {
			return false, err
		}
		if gitDirExists {
			// Ensure the gitdir is empty
			gitDirIsEmpty, err := utils.IsDirEmpty(repo.GitDir)
			if err != nil {
				return false, err
			}
			if !gitDirIsEmpty {
				err = &models.ShitException{Message: fmt.Sprintf("Git directory %s is not empty", repo.GitDir)}
				return false, err
			}
		}
	} else {
		log.Println("Creating shit directory at ", repo.Worktree)
		err = os.MkdirAll(repo.Worktree, 0755)
		if err != nil {
			return false, err
		}
	}

	// Create dirs for branches, objects, refs, refs/heads, refs/tags
	dirs := []string{
		"branches",
		"objects",
		"refs",
		"refs/heads",
		"refs/tags",
	}
	for _, dir := range dirs {
		dirPath := filepath.Join(repo.GitDir, dir)
		err = os.MkdirAll(dirPath, 0755)
		log.Debugf("Creating directory %s", dir)
		if err != nil {
			return false, err
		}
	}

	// Create file description
	descriptionPath := filepath.Join(repo.GitDir, "description")
	descriptionContent := "Unnamed repository; edit this file 'description' to name the repository.\n"
	err = os.WriteFile(descriptionPath, []byte(descriptionContent), 0644)

	// Create file HEAD
	headPath := filepath.Join(repo.GitDir, "HEAD")
	headContent := []byte("ref: refs/heads/main\n")
	err = os.WriteFile(headPath, headContent, 0644)

	// Create the config file
	config := &models.Config{
		Core: models.CoreSection{
			RepositoryFormatVersion: 0,
			FileMode:                false,
			Bare:                    false,
		},
	}

	cfg := ini.Empty()
	err = ini.ReflectFrom(cfg, config)
	if err != nil {
		err = &models.ShitException{Message: fmt.Sprintf("Failed to create config file: %v", err)}
		return false, err
	}
	configFilePath := filepath.Join(repo.GitDir, "config")
	err = cfg.SaveTo(configFilePath)
	if err != nil {
		err := &models.ShitException{Message: fmt.Sprintf("Failed to save config file: %v", err)}
		return false, err
	}
	log.Debugf("Config file saved to %s", configFilePath)
	return true, nil
}
