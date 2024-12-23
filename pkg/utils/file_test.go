package utils

import (
	"os"
	"testing"
)

func TestPathExists(t *testing.T) {
	// Test case: Path exists
	existingFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(existingFile.Name())

	exists, err := PathExists(existingFile.Name())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !exists {
		t.Errorf("Expected file to exist")
	}

	// Test case: Path does not exist
	nonExistentPath := "/non/existent/path"
	exists, err = PathExists(nonExistentPath)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if exists {
		t.Errorf("Expected file to not exist")
	}

	// Test case: Invalid path
	invalidPath := string([]byte{0x00})
	_, err = PathExists(invalidPath)
	if err == nil {
		t.Errorf("Expected an error for invalid path")
	}
}
