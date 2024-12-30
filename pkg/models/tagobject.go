package models

import "fmt"

type ShitTag struct {
	Data        []byte
	TagMetadata *ShitTagMetadata
}

func (b *ShitTag) GetType() string {
	return "tag"
}

func (b *ShitTag) Initialize() error {
	return nil
}

func (b *ShitTag) Serialize(repo *Repository) ([]byte, error) {
	serialized, err := b.TagMetadata.Serialize()
	if err != nil {
		return nil, &ShitException{Message: fmt.Sprintf("Failed to serialize tag object: %v", err)}
	}
	return []byte(serialized), nil
}

func (b *ShitTag) Deserialize(data []byte) error {
	metadata, err := CreateShitTagMetadata(string(data))
	b.TagMetadata = metadata
	if err != nil {
		return &ShitException{Message: fmt.Sprintf("Failed to deserialize tag object: %v", err)}
	}
	return nil
}

func NewShitTag(data []byte) (*ShitTag, error) {
	tagObject := &ShitTag{
		Data: data,
	}
	var err error
	if len(data) == 0 {
		err = tagObject.Initialize()
	} else {
		err = tagObject.Deserialize(data)
	}
	return tagObject, err
}
