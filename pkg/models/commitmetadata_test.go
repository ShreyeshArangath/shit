package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCommitMetadata(t *testing.T) {
	path := "testdata/commit_body"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test resource: %v", err)
	}
	metadata, err := CreateShitCommitMetadata(string(data))
	assert.NoError(t, err)
	assert.NotNil(t, metadata)

	assert.Equal(t, metadata.GetTree(), "29ff16c9c14e2652b22f8b78bb08a5a07930c147")
	assert.Equal(t, metadata.GetParent()[0], "206941306e8a8af65b66eaaaea388a7ae24d49a0")
	assert.Equal(t, metadata.GetAuthor(), "Shreyesh Arangath <shreyesh@shit.dev> 1527025023 +0200")
	assert.Equal(t, metadata.committer, "Shreyesh Arangath <shreyesh@shit.dev> 1527025044 +0200")
	assert.Contains(t, metadata.gpgsignature, "-----BEGIN PGP SIGNATURE-----")
	assert.Contains(t, metadata.gpgsignature, "-----END PGP SIGNATURE-----")
	assert.Equal(t, metadata.message, "Create first draft")
}

func TestSerialize(t *testing.T) {
	metadata := &ShitCommitMetadata{
		tree:         "29ff16c9c14e2652b22f8b78bb08a5a07930c147",
		parent:       []string{"206941306e8a8af65b66eaaaea388a7ae24d49a0"},
		author:       "Shreyesh Arangath <shreyesh@shit.dev> 1527025023 +0200",
		committer:    "Shreyesh Arangath <shreyesh@shit.dev> 1527025044 +0200",
		gpgsignature: "-----BEGIN PGP SIGNATURE-----\n\niQIzBAABCAAdFiEExwXquOM8bWb4Q2zVGxM2FxoLkGQFAlsEjZQACgkQGxM2FxoL\nkGQdcBAAqPP+ln4nGDd2gETXjvOpOxLzIMEw4A9gU6CzWzm+oB8mEIKyaH0UFIPh\nrNUZ1j7/ZGFNeBDtT55LPdPIQw4KKlcf6kC8MPWP3qSu3xHqx12C5zyai2duFZUU\nwqOt9iCFCscFQYqKs3xsHI+ncQb+PGjVZA8+jPw7nrPIkeSXQV2aZb1E68wa2YIL\n3eYgTUKz34cB6tAq9YwHnZpyPx8UJCZGkshpJmgtZ3mCbtQaO17LoihnqPn4UOMr\nV75R/7FjSuPLS8NaZF4wfi52btXMSxO/u7GuoJkzJscP3p4qtwe6Rl9dc1XC8P7k\nNIbGZ5Yg5cEPcfmhgXFOhQZkD0yxcJqBUcoFpnp2vu5XJl2E5I/quIyVxUXi6O6c\n/obspcvace4wy8uO0bdVhc4nJ+Rla4InVSJaUaBeiHTW8kReSFYyMmDCzLjGIu1q\ndoU61OM3Zv1ptsLu3gUE6GU27iWYj2RWN3e3HE4Sbd89IFwLXNdSuM0ifDLZk7AQ\nWBhRhipCCgZhkj9g2NEk7jRVslti1NdN5zoQLaJNqSwO1MtxTmJ15Ksk3QP6kfLB\nQ52UWybBzpaP9HEd4XnR+HuQ4k2K0ns2KgNImsNvIyFwbpMUyUWLMPimaV1DWUXo\n5SBjDB/V/W2JBFR+XKHFJeFwYhj7DD/ocsGr4ZMx/lgc8rjIBkI=\n=lgTX\n-----END PGP SIGNATURE-----\n\n",
		message:      "Create first draft",
	}
	bytes, err := metadata.Serialize()
	assert.NoError(t, err)
	assert.NotNil(t, bytes)
	serialized_metadata_str := string(bytes)

	path := "testdata/commit_body"
	data, err := os.ReadFile(path)
	assert.Equal(t, serialized_metadata_str, string(data))
}

func TestSerDe(t *testing.T) {
	path := "testdata/commit_body"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test resource: %v", err)
	}
	metadata, err := CreateShitCommitMetadata(string(data))
	assert.NoError(t, err)
	assert.NotNil(t, metadata)

	bytes, err := metadata.Serialize()
	// os.WriteFile("testdata/serialized_commit_body", bytes, 0644)
	assert.NoError(t, err)
	assert.NotNil(t, bytes)
	serialized_metadata_str := string(bytes)
	assert.Equal(t, serialized_metadata_str, string(data))
}

func TestDeSer(t *testing.T) {
	path := "testdata/commit_body"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test resource: %v", err)
	}
	metadata, err := CreateShitCommitMetadata(string(data))
	assert.NoError(t, err)
	assert.NotNil(t, metadata)

	bytes, err := metadata.Serialize()
	assert.NoError(t, err)
	assert.NotNil(t, bytes)
	serialized_metadata_str := string(bytes)
	assert.Equal(t, serialized_metadata_str, string(data))

	metadata2, err := CreateShitCommitMetadata(serialized_metadata_str)
	assert.NoError(t, err)
	assert.NotNil(t, metadata2)
	assert.Equal(t, metadata, metadata2)
}
