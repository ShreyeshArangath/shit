package models

import (
	"fmt"
)

type ShitCommit struct {
	Data           []byte
	CommitMetadata *ShitCommitMetadata
}

func (b *ShitCommit) GetType() string {
	return "commit"
}

func (b *ShitCommit) Initialize() error {
	return nil
}

func (b *ShitCommit) Serialize(repo *Repository) ([]byte, error) {
	serialized, err := b.CommitMetadata.Serialize()
	if err != nil {
		return nil, &ShitException{Message: fmt.Sprintf("Failed to serialize commit object: %v", err)}
	}
	return []byte(serialized), nil
}

func (b *ShitCommit) Deserialize(data []byte) error {
	metadata, err := CreateShitCommitMetadata(string(data))
	b.CommitMetadata = metadata
	if err != nil {
		return &ShitException{Message: fmt.Sprintf("Failed to deserialize commit object: %v", err)}
	}
	return nil
}

func NewShitCommit(data []byte) (*ShitCommit, error) {
	commitObject := &ShitCommit{
		Data: data,
	}
	var err error
	if len(data) == 0 {
		err = commitObject.Initialize()
	} else {
		err = commitObject.Deserialize(data)
	}
	return commitObject, err
}
