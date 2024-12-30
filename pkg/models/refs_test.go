package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ShreyeshArangath/shit/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/ini.v1"
)

func setup(t *testing.T) string {
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

	// copy the refs dir from the testdata dir to the temp dir
	refsDir := filepath.Join("testdata", "refs")
	err = utils.CopyDir(refsDir, filepath.Join(gitDir, "refs"))
	return tempDir
}

func TestResolveRefHappyPath(t *testing.T) {
	tempDir := setup(t)
	repo, err := CreateRepository(tempDir, false)
	assert.NoError(t, err)
	hash, err := ResolveRef(repo, filepath.Join(tempDir, ".git/refs/remote/origin/main"))
	assert.NoError(t, err)
	assert.Equal(t, "0fa6863b05513c03ea85033be9cf95d2ca035e27", hash)
}
