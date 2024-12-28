package models

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

type TreeLeaf struct {
	Mode string
	Path string
	Sha  string
}

func ParseLeaf(raw []byte, start int) (int, TreeLeaf, error) {
	x := bytes.IndexByte(raw[start:], ' ')
	if x == -1 {
		return 0, TreeLeaf{}, fmt.Errorf("space not found for mode")
	}
	x += start

	// Ensure mode length is valid
	if x-start != 5 && x-start != 6 {
		return 0, TreeLeaf{}, fmt.Errorf("invalid mode length: %d", x-start)
	}

	// Read the mode and normalize to six bytes if needed
	mode := raw[start:x]
	if len(mode) == 5 {
		mode = append([]byte(" "), mode...)
	}

	// Find the NULL terminator of the path
	y := bytes.IndexByte(raw[x+1:], '\x00')
	if y == -1 {
		return 0, TreeLeaf{}, fmt.Errorf("NULL terminator not found for path")
	}
	y += x + 1

	// Read the path
	path := raw[x+1 : y]

	// Ensure sufficient bytes for SHA
	if y+21 > len(raw) {
		return 0, TreeLeaf{}, fmt.Errorf("insufficient bytes for SHA")
	}

	// Read the SHA
	rawSha := raw[y+1 : y+21]
	sha := hex.EncodeToString(rawSha)

	return y + 21, TreeLeaf{
		Mode: string(mode),
		Path: string(path),
		Sha:  sha,
	}, nil

}

func (t TreeLeaf) Serialize() ([]byte, error) {
	// mode + space + path + null + sha
	var buffer bytes.Buffer

	// Add the mode
	buffer.WriteString(t.Mode)
	buffer.WriteByte(' ') // Space separator

	// Add the path
	buffer.WriteString(t.Path)
	buffer.WriteByte(0) // NULL terminator

	// Convert the SHA from hexadecimal to bytes
	shaBytes := make([]byte, 20)
	_, err := fmt.Sscanf(t.Sha, "%40x", &shaBytes)
	if err != nil {
		return nil, fmt.Errorf("invalid SHA: %w", err)
	}
	buffer.Write(shaBytes)
	return buffer.Bytes(), nil
}

func (t TreeLeaf) SortKey() string {
	if t.Mode[:2] == "10" {
		return t.Path
	}
	return t.Path + "/"
}
