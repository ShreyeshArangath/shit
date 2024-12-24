package models

type ShitBlob struct {
	Data []byte
}

func (b *ShitBlob) GetType() string {
	return "blob"
}

func (b *ShitBlob) Initialize() error {
	// No initialization needed for a blob
	return nil
}

func (b *ShitBlob) Serialize(repo *Repository) ([]byte, error) {
	return b.Data, nil
}

func (b *ShitBlob) Deserialize(data []byte) error {
	b.Data = data
	return nil
}

func NewShitBlob(data []byte) (*ShitBlob, error) {
	blob := &ShitBlob{
		Data: data,
	}
	var err error
	if len(data) == 0 {
		err = blob.Initialize()
	} else {
		err = blob.Deserialize(data)
	}
	return blob, err
}
