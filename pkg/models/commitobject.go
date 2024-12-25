package models

import (
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type ShitCommit struct {
	Data []byte
	KVLM *KVLM
}

func (b *ShitCommit) GetType() string {
	return "commit"
}

func (b *ShitCommit) Initialize() error {
	b.KVLM = CreateKVLM(orderedmap.New[string, []string]())
	return nil
}

func (b *ShitCommit) Serialize(repo *Repository) ([]byte, error) {
	serialized, err := KVLMSerialize(b.KVLM)
	if err != nil {
		return nil, err
	}
	return []byte(serialized), nil
}

func (b *ShitCommit) Deserialize(data []byte) error {
	kvlm, err := KVLMDeserialize(string(data))
	if err != nil {
		return err
	}
	b.KVLM = kvlm
	return nil
}
