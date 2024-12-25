package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVLMDeserializeHappyPath(t *testing.T) {
	path := "testdata/commit_body"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test resource: %v", err)
	}
	kvlm, err := KVLMDeserialize(string(data))
	assert.NoError(t, err)
	assert.NotNil(t, kvlm)
	assert.Equal(t, kvlm.OrderedMap.Len(), 6)

	treeVal, _ := kvlm.OrderedMap.Get("tree")
	assert.Equal(t, treeVal[0], "29ff16c9c14e2652b22f8b78bb08a5a07930c147")

	parent, _ := kvlm.OrderedMap.Get("parent")
	assert.Equal(t, parent[0], "206941306e8a8af65b66eaaaea388a7ae24d49a0")

	author, _ := kvlm.OrderedMap.Get("author")
	assert.Equal(t, author[0], "Shreyesh Arangath <shreyesh@shit.dev> 1527025023 +0200")

	committer, _ := kvlm.OrderedMap.Get("committer")
	assert.Equal(t, committer[0], "Shreyesh Arangath <shreyesh@shit.dev> 1527025044 +0200")

	gpgsig, _ := kvlm.OrderedMap.Get("gpgsig")
	assert.Contains(t, gpgsig[0], "-----BEGIN PGP SIGNATURE-----")
	assert.Contains(t, gpgsig[0], "-----END PGP SIGNATURE-----")

	message, _ := kvlm.OrderedMap.Get("")
	assert.Equal(t, message[0], "Create first draft")
}
