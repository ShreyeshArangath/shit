package models

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ShreyeshArangath/shit/pkg/utils"
)

const (
	REFERENCE_PREFIX = "ref: "
)

type RefMap map[string]interface{} // Allows for nested maps

// ResolveRef resolves a Git reference to its corresponding commit hash.
// It takes a Repository object and a reference string as input, and returns
// the resolved commit hash as a string, or an error if the resolution fails.
//
// The function first computes the relative path of the reference file within
// the repository's Git directory. It then reads the content of the reference
// file, which should contain the commit hash or another reference.
//
// If the reference points to another reference (indicated by the REFERENCE_PREFIX),
// the function recursively resolves the new reference.
//
// Parameters:
// - repo: A pointer to the Repository object.
// - ref: The reference string to resolve.
//
// Returns:
// - A string containing the resolved commit hash.
// - An error if the resolution process fails.
func ResolveRef(repo *Repository, ref string) (string, error) {
	relPath, err := filepath.Rel(repo.GitDir, ref)
	path, err := repo.RepoFile(false, relPath)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}
	// Drop the last \n
	data = data[:len(data)-1]
	datastr := string(data)
	if strings.HasPrefix(datastr, REFERENCE_PREFIX) {
		newRef := filepath.Join(repo.GitDir, datastr[len(REFERENCE_PREFIX):])
		return ResolveRef(repo, newRef)
	}
	return datastr, nil
}

// ListRef lists all references in the given repository and path.
// If the path is empty, it defaults to the "refs" directory in the repository.
//
// Parameters:
//   - repo: A pointer to the Repository struct.
//   - path: A string representing the path to list references from.
//
// Returns:
//   - RefMap: A map of references found in the specified path.
//   - error: An error if any occurred during the operation.
//
// The function reads the directory contents, sorts them alphabetically,
// and iterates through each entry. If an entry is a directory, it recurses
// into the subdirectory. If an entry is a file, it resolves the reference
// using the ResolveRef helper function.
func ListRef(repo *Repository, path string) (RefMap, error) {
	if path == "" {
		var err error
		path, err = repo.RepoFile(false, "refs")
		if err != nil {
			return nil, err
		}
	}

	ret := make(RefMap)

	// Read directory contents and sort them
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Sort entries alphabetically (mimics sorted(os.listdir))
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		if ok, _ := utils.PathExists(fullPath); !ok {
			continue
		}
		if entry.IsDir() {
			// Recurse into subdirectories
			subRefs, err := ListRef(repo, fullPath)
			if err != nil {
				return nil, err
			}
			ret[entry.Name()] = subRefs
		} else {
			// Resolve the reference
			ref, err := ResolveRef(repo, fullPath) // Assume ResolveRef is a helper function
			if err != nil {
				return nil, err
			}
			ret[entry.Name()] = ref
		}
	}

	return ret, nil
}

// CreateRef creates a new reference in the given repository.
// It takes a repository pointer, a reference name, and a SHA string as inputs.
// The function constructs the reference path, writes the SHA to the reference file,
// and returns an error if any operation fails.
//
// Parameters:
//   - repo: A pointer to the Repository where the reference will be created.
//   - ref: The name of the reference to be created.
//   - sha: The SHA string to be written to the reference file.
//
// Returns:
//   - error: An error if any operation fails, otherwise nil.
func CreateRef(repo *Repository, ref string, sha string) error {
	refName := filepath.Join("refs/", ref)
	refPath, err := repo.RepoFile(true, refName)
	if err != nil {
		return err
	}
	data := fmt.Sprintf("%s\n", sha)
	err = os.WriteFile(refPath, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}
