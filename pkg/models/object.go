package models

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ShreyeshArangath/shit/pkg/utils"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

type Object interface {
	GetType() string
	Initialize() error
	Serialize(repo *Repository) ([]byte, error)
	Deserialize(data []byte) error
}

// ObjectFactory creates an object of the specified type from the provided data.
// It returns an Object and an error if the object type is unknown or if there is an issue
// with the creation of the object.
//
// Parameters:
//   - objectType: A string representing the type of object to create. Valid values are "commit", "tree", "tag", and "blob".
//   - data: A byte slice containing the data to be used for creating the object.
//
// Returns:
//   - Object: The created object that implements the Object interface.
//   - error: An error if the object type is unknown or if there is an issue with the creation of the object.
func ObjectFactory(objectType string, data []byte) (Object, error) {
	// To implement the `Object` interface and its methods, we need to define the different types of objects (commit, tree, tag, blob) and their respective serialization and deserialization methods. Below is a basic implementation of the `Object` interface and its methods for each type of object.
	switch objectType {
	case "commit":
		return NewShitCommit(data)
	case "tree":
		return NewShitTree(data)
	case "tag":
		return nil, nil
	case "blob":
		return NewShitBlob(data)
	default:
		return nil, &ShitException{Message: fmt.Sprintf("Unknown type %s", objectType)}
	}
}

// ObjectRead reads an object from the repository using its SHA-1 hash.
// It first constructs the path to the object file and checks if the file exists.
// If the file exists, it reads the binary data and decompresses it using zlib.
// The decompressed data is then parsed to extract the object type and data.
// If the object is malformed or any error occurs during the process, an error is returned.
//
// Parameters:
//   - repo: A pointer to the Repository from which the object is to be read.
//   - sha: The SHA-1 hash of the object to be read.
//
// Returns:
//   - Object: The object read from the repository.
//   - error: An error if any occurs during the reading process.
func ObjectRead(repo *Repository, sha string) (Object, error) {
	path := repo.GetRepoPath("objects", sha[:2], sha[2:])
	isFile, err := utils.IsFile(path)
	if err != nil {
		return nil, err
	}
	if !isFile {
		return nil, nil
	}

	data, err := readBinaryFile(path)
	if err != nil {
		return nil, err
	}

	var decompressedData bytes.Buffer
	decompressedData, err = zlibDecompress(data)
	if err != nil {
		return nil, err
	}

	decompressedDataStr := decompressedData.String()
	spaceIndex := strings.Index(decompressedDataStr, " ")
	objectType := decompressedDataStr[:spaceIndex]
	if spaceIndex == -1 {
		return nil, &ShitException{Message: "Malformed object missing space"}
	}

	// Read and validate the object size
	objectSizeIndex := strings.Index(decompressedDataStr, "\x00")
	_, err = strconv.Atoi(decompressedDataStr[spaceIndex+1 : objectSizeIndex])
	if err != nil {
		return nil, err
	}
	objectData := decompressedDataStr[objectSizeIndex+1:]
	return ObjectFactory(objectType, []byte(objectData))
}

// ObjectWrite writes an object to the repository.
// It serializes the object, computes its hash, and stores it in the repository if it doesn't already exist.
//
// Parameters:
//   - object: The object to be written.
//   - repo: The repository where the object will be stored.
//
// Returns:
//   - string: The SHA-1 hash of the object.
//   - error: An error if the operation fails.
func ObjectWrite(object Object, repo *Repository) (string, error) {
	data, err := object.Serialize(repo)
	if err != nil {
		return "", err
	}

	// Add the header
	header := fmt.Sprintf("%s %s\x00", object.GetType(), strconv.Itoa(len(data)))
	// Add the contents of the object to the header
	contents := append([]byte(header), data...)
	// Compute the hash
	sha1 := sha1.New()
	sha1.Write(contents)
	sha := hex.EncodeToString(sha1.Sum(nil))
	objectPath, err := repo.RepoFile(true, "objects", sha[:2], sha[2:])
	if err != nil {
		return "", err
	}

	objectPathExists, err := utils.PathExists(objectPath)
	if err != nil {
		return "", err
	}

	if !objectPathExists {
		// Write the object to the store
		compressedData, err := zlibCompress(contents)
		if err != nil {
			return "", err
		}
		err = writeBinaryFile(objectPath, compressedData)
		if err != nil {
			return "", err
		}
	}
	return sha, nil
}

func ObjectFind(repo *Repository, name string, objecttype string, follow bool) (string, error) {
	return name, nil
}

// ObjectHash computes the hash of an object in the repository.
//
// Parameters:
//   - repo: A pointer to the Repository where the object resides.
//   - objectype: A string representing the type of the object.
//   - path: A string representing the file path to read the object data from.
//
// Returns:
//   - A string representing the hash of the object.
//   - An error if any occurs during the process of reading the file, creating the object, or writing the object.
func ObjectHash(repo *Repository, objectype string, path string) (string, error) {
	data, err := readBinaryFile(path)
	object, err := ObjectFactory(objectype, data)
	if err != nil {
		return "", err
	}
	return ObjectWrite(object, repo)
}

func readBinaryFile(path string) ([]byte, error) {
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

func writeBinaryFile(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(data)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

// Decompresses the object data using zlib and returns the decompressed data.
func zlibDecompress(data []byte) (bytes.Buffer, error) {
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

// zlibCompress compresses the given data using zlib.
func zlibCompress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
