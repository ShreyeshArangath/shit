package models

import (
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
		return ResolveRef(repo, datastr[len(REFERENCE_PREFIX):])
	}
	return datastr, nil
}

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
