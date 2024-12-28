package models

type ShitTree struct {
	Data  []byte
	Items []TreeLeaf
}

func (t *ShitTree) GetType() string {
	return "tree"
}

func (t *ShitTree) Initialize() error {
	return nil
}

func (t *ShitTree) Serialize(repo *Repository) ([]byte, error) {
	return nil, nil
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
