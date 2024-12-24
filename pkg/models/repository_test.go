package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ShreyeshArangath/shit/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
)

func setupTest(t *testing.T) (string, string) {
	// Setup
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, utils.GIT_DIR_NAME)
	err := os.Mkdir(gitDir, 0755)
	assert.NoError(t, err)

	// Create a dummy config file
	configFilePath := filepath.Join(gitDir, utils.CONFIG_FILE_NAME)
	cfg := ini.Empty()
	cfg.Section("core").Key("repositoryformatversion").SetValue("0")
	err = cfg.SaveTo(configFilePath)
	assert.NoError(t, err)
	return tempDir, gitDir
}

func TestCreateRepositoryHappyPath(t *testing.T) {
	tempDir, gitDir := setupTest(t)
	repo, err := CreateRepository(tempDir, false)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, tempDir, repo.Worktree)
	assert.Equal(t, gitDir, repo.GitDir)
}

func TestCreateRepositoryInvalidPath(t *testing.T) {
	_, err := CreateRepository("/invalid/path", false)
	assert.Error(t, err)
}

func TestCreateRepositoryNonExistentConfigFile(t *testing.T) {
	_, err := CreateRepository(t.TempDir(), false)
	assert.Error(t, err)
}

func TestCreateRepositoryUnsupportedFormatVersion(t *testing.T) {
	tempDir, gitDir := setupTest(t)
	cfg, err := ini.Load(filepath.Join(gitDir, utils.CONFIG_FILE_NAME))
	assert.NoError(t, err)
	cfg.Section("core").Key("repositoryformatversion").SetValue("1")
	err = cfg.SaveTo(filepath.Join(gitDir, utils.CONFIG_FILE_NAME))
	assert.NoError(t, err)
	_, err = CreateRepository(tempDir, false)
	assert.Error(t, err)
}

func TestRepoDirExists(t *testing.T) {
	tempDir, gitDir := setupTest(t)
	repo, err := CreateRepository(tempDir, false)

	gitDirNewDir := filepath.Join(gitDir, "existing")
	err = os.Mkdir(gitDirNewDir, 0755)
	assert.NoError(t, err)

	dirPath, err := repo.RepoDir(false, "existing")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(gitDir, "existing"), dirPath)
}

func TestRepoDirDNECreate(t *testing.T) {
	tempDir, gitDir := setupTest(t)
	repo, err := CreateRepository(tempDir, false)

	dirPath, err := repo.RepoDir(true, "newdir")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(gitDir, "newdir"), dirPath)
}

func TestRepoDirDNEDoNotCreate(t *testing.T) {
	tempDir, _ := setupTest(t)
	repo, err := CreateRepository(tempDir, false)

	_, err = repo.RepoDir(false, "nonexistent")
	assert.Error(t, err)
}

func TestRepoFileExists(t *testing.T) {
	tempDir, gitDir := setupTest(t)
	repo, err := CreateRepository(tempDir, false)
	gitDirNewDir := filepath.Join(gitDir, "existing")
	err = os.Mkdir(gitDirNewDir, 0755)
	filePath, err := repo.RepoFile(false, "existing", "existingfile")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(gitDir, "existing", "existingfile"), filePath)
}

func TestRepoFileDNECreate(t *testing.T) {
	tempDir, gitDir := setupTest(t)
	repo, err := CreateRepository(tempDir, false)
	filePath, err := repo.RepoFile(true, "helloworld", "newfile")
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(gitDir, "helloworld", "newfile"), filePath)
}

func TestRepoFileDNEDoNotCreate(t *testing.T) {
	tempDir, _ := setupTest(t)
	repo, err := CreateRepository(tempDir, false)
	_, err = repo.RepoFile(false, "helloworld", "nonexistentfile")
	assert.Error(t, err)
}

func TestRepoFindHappyPath(t *testing.T) {
	tempDir, gitDir := setupTest(t)
	repo, err := RepoFind(tempDir, true)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, tempDir, repo.Worktree)
	assert.Equal(t, gitDir, repo.GitDir)
}

func TestRepoFindNonExistentPath(t *testing.T) {
	_, err := RepoFind("/invalid/path", true)
	assert.Error(t, err)
}

func TestRepoFindNoGitDir(t *testing.T) {
	tempDir := t.TempDir()
	_, err := RepoFind(tempDir, true)
	assert.Error(t, err)
}

func TestRepoFindNoGitDirNotRequired(t *testing.T) {
	tempDir := t.TempDir()
	repo, err := RepoFind(tempDir, false)
	assert.NoError(t, err)
	assert.Nil(t, repo)
}

func TestRepoFindInParentDir(t *testing.T) {
	tempDir, gitDir := setupTest(t)
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	assert.NoError(t, err)

	repo, err := RepoFind(subDir, true)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, tempDir, repo.Worktree)
	assert.Equal(t, gitDir, repo.GitDir)
}
