package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShitMetadataCreate(t *testing.T) {
	data, err := os.ReadFile("testdata/tag")
	assert.NoError(t, err)
	metadata, err := CreateShitTagMetadata(string(data))
	assert.NoError(t, err)
	assert.Equal(t, "1519c056ada87c66a34004a868c4f1cb188fade6", metadata.object)
	assert.Equal(t, "commit", metadata.objecttype)
	assert.Equal(t, "v1.0", metadata.tag)
	assert.Equal(t, "Shreyesh Arangath <sarangath@google.com> 1735610175 +0300", metadata.tagger)
	assert.Equal(t, "test", metadata.message)
}

func TestShitMetadataSerialize(t *testing.T) {
	data, err := os.ReadFile("testdata/tag")
	assert.NoError(t, err)
	metadata, err := CreateShitTagMetadata(string(data))
	assert.NoError(t, err)
	serialized, err := metadata.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, string(data), serialized)
}
