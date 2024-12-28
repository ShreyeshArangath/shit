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
	// file mode terminated by  0x20
	// path terminated by null
	// sha 1  40 bytes
	// Find the space terminator of the mode
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
	return nil, nil
}
