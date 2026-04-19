package core

type Blob struct {
	BaseObject
}

func NewBlob(content []byte) *Blob {
	return &Blob{
		BaseObject: BaseObject{
			Type:    "blob",
			Content: content,
		},
	}
}
