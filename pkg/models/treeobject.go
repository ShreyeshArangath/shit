package models

import "sort"

type ShitTree struct {
	Data  []byte
	Items []TreeLeaf
}

func (t *ShitTree) GetType() string {
	return "tree"
}

func (t *ShitTree) Initialize() error {
	t.Items = make([]TreeLeaf, 0)
	return nil
}

func (t *ShitTree) Serialize(repo *Repository) ([]byte, error) {
	t.sortTreeLeaf()
	output := make([]byte, 0)
	for _, item := range t.Items {
		serializedLeaf, err := item.Serialize()
		if err != nil {
			return nil, err
		}
		output = append(output, serializedLeaf...)
	}
	return output, nil
}

func (t *ShitTree) sortTreeLeaf() {
	sort.Slice(t.Items, func(i, j int) bool {
		return t.Items[i].SortKey() < t.Items[j].SortKey()
	})
}

func (t *ShitTree) Deserialize(data []byte) error {
	pos := 0
	var items []TreeLeaf
	for pos < len(data) {
		index, leaf, err := ParseLeaf(data, pos)
		pos = index
		if err != nil {
			return err
		}
		items = append(items, leaf)
	}
	t.Items = items
	return nil
}
