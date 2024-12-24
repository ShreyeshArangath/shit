package models

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ShreyeshArangath/shit/pkg/utils"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Object interface {
	initialize()
	serialize(repo *Repository) []byte
	deserialize(data []byte)
}

func ObjectFactory(objectType string, data []byte) (*Object, error) {
	// To implement the `Object` interface and its methods, we need to define the different types of objects (commit, tree, tag, blob) and their respective serialization and deserialization methods. Below is a basic implementation of the `Object` interface and its methods for each type of object.
	switch objectType {
	case "commit":
		return nil, nil
	case "tree":
		return nil, nil
	case "tag":
		return nil, nil
	case "blob":
		return nil, nil
	default:
		return nil, &ShitException{Message: fmt.Sprintf("Unknown type %s", objectType)}
	}
}

// Reads the object SHA and returns an Object representation of it.
func ObjectRead(repo *Repository, sha string) (*Object, error) {
	path := repo.GetRepoPath("objects", sha[:2], sha[2:])
	isFile, err := utils.IsFile(path)
	if err != nil {
		return nil, err
	}
	if !isFile {
		return nil, nil
	}

	data, err := ReadBinaryFile(path)
	if err != nil {
		return nil, err
	}

	var decompressedData bytes.Buffer
	decompressedData, err = ZlibDecompress(data)
	if err != nil {
		return nil, err
	}

	decompressedDataStr := decompressedData.String()
	// Read the object type
	spaceIndex := strings.Index(decompressedDataStr, " ")
	objectType := decompressedDataStr[:spaceIndex]
	if spaceIndex == -1 {
		return nil, &ShitException{Message: "Malformed object missing space"}
	}

	// Read and validate the object size
	objectSizeIndex := strings.Index(decompressedDataStr, "\x00")
	objectSize, err := strconv.Atoi(decompressedDataStr[spaceIndex+1 : objectSizeIndex])
	if err != nil {
		return nil, err
	}
	if objectSize != (len(decompressedDataStr) - objectSizeIndex) {
		return nil, &ShitException{Message: fmt.Sprintf("Malformed object (bad length) %s", sha)}
	}

	objectData := decompressedDataStr[objectSizeIndex+1:]
	return ObjectFactory(objectType, []byte(objectData))
}

func ReadBinaryFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var buffer []byte
	reader := bufio.NewReader(file)
	buf := make([]byte, 1024) // Read in chunks of 1 KB
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		buffer = append(buffer, buf[:n]...)
	}
	return buffer, nil
}

// Decompresses the object data using zlib and returns the decompressed data.
func ZlibDecompress(data []byte) (bytes.Buffer, error) {
	reader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return bytes.Buffer{}, err
	}
	defer reader.Close()

	var out bytes.Buffer
	_, err = out.ReadFrom(reader)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return out, nil
}
