package core

type Tree struct {
	BaseObject
}

func NewTree() *Tree {
	return &Tree{
		BaseObject: BaseObject{
			Type: "tree",
		},
	}
}
