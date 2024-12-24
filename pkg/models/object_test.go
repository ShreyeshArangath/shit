package models

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ShreyeshArangath/shit/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func setupObjectTest(t *testing.T) (string, string, *Repository) {
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, utils.GIT_DIR_NAME)
	err := os.Mkdir(gitDir, 0755)
	assert.NoError(t, err)

	repo, err := CreateRepository(tempDir, true)
	assert.NoError(t, err)
	return tempDir, gitDir, repo
}

// TODO: Update this test case to use the actual implementation of the tree object
func TestObjectFactoryCommit(t *testing.T) {
	_, err := ObjectFactory("commit", []byte{})
	assert.NoError(t, err)
}

// TODO: Update this test case to use the actual implementation of the tree object
func TestObjectFactoryTree(t *testing.T) {
	_, err := ObjectFactory("tree", []byte{})
	assert.NoError(t, err)
}

// TODO: Update this test case to use the actual implementation of the tag object
func TestObjectFactoryTag(t *testing.T) {
	_, err := ObjectFactory("tag", []byte{})
	assert.NoError(t, err)
}

// TODO: Update this test case to use the actual implementation of the blob object
func TestObjectFactoryBlob(t *testing.T) {
	_, err := ObjectFactory("blob", []byte{})
	assert.NoError(t, err)
}

func TestObjectFactoryInvalid(t *testing.T) {
	_, err := ObjectFactory("invalid", []byte{})
	assert.Error(t, err)
}

func TestObjectReadValidSHA(t *testing.T) {
	_, gitDir, repo := setupObjectTest(t)

	// Create a valid object file
	sha := "1234567890abcdef1234567890abcdef12345678"
	objectPath := filepath.Join(gitDir, "objects", sha[:2], sha[2:])
	err := os.MkdirAll(filepath.Dir(objectPath), 0755)
	assert.NoError(t, err)

	data := []byte("commit\x2010\x00test data")
	compressedData := new(bytes.Buffer)
	writer := zlib.NewWriter(compressedData)
	_, err = writer.Write(data)
	assert.NoError(t, err)
	writer.Close()
	err = os.WriteFile(objectPath, compressedData.Bytes(), 0644)
	assert.NoError(t, err)

	_, err = ObjectRead(repo, sha)
	assert.NoError(t, err)
	// assert.NotNil(t, obj)
}

func TestObjectReadInvalidSHA(t *testing.T) {
	_, _, repo := setupObjectTest(t)
	_, err := ObjectRead(repo, "invalidsha")
	assert.Error(t, err)
}

func TestObjectReadMalformedData(t *testing.T) {
	_, gitDir, repo := setupObjectTest(t)

	// Create a malformed object file
	sha := "1234567890abcdef1234567890abcdef12345678"
	objectPath := filepath.Join(gitDir, "objects", sha[:2], sha[2:])
	err := os.MkdirAll(filepath.Dir(objectPath), 0755)
	assert.NoError(t, err)

	err = os.WriteFile(objectPath, []byte("malformed data"), 0644)
	assert.NoError(t, err)

	_, err = ObjectRead(repo, sha)
	assert.Error(t, err)
}

func TestReadBinaryFileExists(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile")
	err := os.WriteFile(filePath, []byte("test data"), 0644)
	assert.NoError(t, err)

	data, err := readBinaryFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test data"), data)
}

func TestReadBinaryFileNotExists(t *testing.T) {
	tempDir := t.TempDir()
	_, err := readBinaryFile(filepath.Join(tempDir, "nonexistent"))
	assert.Error(t, err)
}

func TestZlibDecompressValidData(t *testing.T) {
	data := []byte("test data")
	var compressedData bytes.Buffer
	writer := zlib.NewWriter(&compressedData)
	_, err := writer.Write(data)
	assert.NoError(t, err)
	writer.Close()

	decompressedData, err := zlibDecompress(compressedData.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, data, decompressedData.Bytes())
}

func TestZlibDecompressInvalidData(t *testing.T) {
	_, err := zlibDecompress([]byte("invalid data"))
	assert.Error(t, err)
}

func TestZlibCompressValidData(t *testing.T) {
	data := []byte("test data")
	compressedData, err := zlibCompress(data)
	assert.NoError(t, err)

	// Decompress to verify
	decompressedData, err := zlibDecompress(compressedData)
	assert.NoError(t, err)
	assert.Equal(t, data, decompressedData.Bytes())
}

func TestZlibCompressEmptyData(t *testing.T) {
	data := []byte("")
	compressedData, err := zlibCompress(data)
	assert.NoError(t, err)

	// Decompress to verify
	decompressedData, err := zlibDecompress(compressedData)
	assert.NoError(t, err)
	assert.Equal(t, data, decompressedData.Bytes())
}

func TestWriteBinaryFileValid(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "testfile")
	data := []byte("test data")

	err := writeBinaryFile(filePath, data)
	assert.NoError(t, err)

	readData, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, data, readData)
}

func TestWriteBinaryFileInvalidPath(t *testing.T) {
	invalidPath := "/invalid/path/testfile"
	data := []byte("test data")

	err := writeBinaryFile(invalidPath, data)
	assert.Error(t, err)
}

func TestObjectWriteHappyPath(t *testing.T) {
	_, gitDir, repo := setupObjectTest(t)

	// Create a fake object
	fakeObject := &FakeObject{
		Type: "blob",
		Data: []byte("test data"),
	}

	sha, err := ObjectWrite(fakeObject, repo)
	assert.NoError(t, err)
	assert.NotEmpty(t, sha)

	// Verify the object was written to the correct path
	objectPath := filepath.Join(gitDir, "objects", sha[:2], sha[2:])
	objectPathExists, err := utils.PathExists(objectPath)
	assert.NoError(t, err)
	assert.True(t, objectPathExists)

	// Verify the contents of the written object
	compressedData, err := os.ReadFile(objectPath)
	assert.NoError(t, err)

	decompressedData, err := zlibDecompress(compressedData)
	assert.NoError(t, err)

	expectedHeader := fmt.Sprintf("blob %d\x00", len(fakeObject.Data))
	expectedContents := append([]byte(expectedHeader), fakeObject.Data...)
	assert.Equal(t, expectedContents, decompressedData.Bytes())
}

func TestObjectHashValid(t *testing.T) {
	_, gitDir, repo := setupObjectTest(t)

	// Create a temporary file with test data
	tempFile := filepath.Join(gitDir, "testfile")
	data := []byte("test data")
	err := os.WriteFile(tempFile, data, 0644)
	assert.NoError(t, err)

	// Call ObjectHash with the temporary file
	sha, err := ObjectHash(repo, "blob", tempFile)
	assert.NoError(t, err)
	assert.NotEmpty(t, sha)

	// Verify the object was written to the correct path
	objectPath := filepath.Join(gitDir, "objects", sha[:2], sha[2:])
	objectPathExists, err := utils.PathExists(objectPath)
	assert.NoError(t, err)
	assert.True(t, objectPathExists)

	// Verify the contents of the written object
	compressedData, err := os.ReadFile(objectPath)
	assert.NoError(t, err)

	decompressedData, err := zlibDecompress(compressedData)
	assert.NoError(t, err)

	expectedHeader := fmt.Sprintf("blob %d\x00", len(data))
	expectedContents := append([]byte(expectedHeader), data...)
	assert.Equal(t, expectedContents, decompressedData.Bytes())
}

func TestObjectHashInvalidObjectType(t *testing.T) {
	_, gitDir, repo := setupObjectTest(t)

	// Create a temporary file with test data
	tempFile := filepath.Join(gitDir, "testfile")
	data := []byte("test data")
	err := os.WriteFile(tempFile, data, 0644)
	assert.NoError(t, err)

	// Call ObjectHash with an invalid object type
	_, err = ObjectHash(repo, "invalidtype", tempFile)
	assert.Error(t, err)
}

// TODO: Update this test case to use the actual implementation of the blob object
type FakeObject struct {
	Type string
	Data []byte
}

func (o *FakeObject) Serialize(repo *Repository) ([]byte, error) {
	return o.Data, nil
}

func (o *FakeObject) Deserialize(data []byte) error {
	o.Data = data
	return nil
}

func (o *FakeObject) GetType() string {
	return o.Type
}

func (o *FakeObject) Initialize() error { return nil }
