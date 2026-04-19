package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type GitObject interface {
	GetType() string
	GetContent() []byte
	Serialize() ([]byte, error)
	Hash() string
}

type BaseObject struct {
	Type    string
	Content []byte
}

func (b *BaseObject) GetType() string    { return b.Type }
func (b *BaseObject) GetContent() []byte { return b.Content }

func (b *BaseObject) Serialize() ([]byte, error) {
	header := fmt.Sprintf("%s %d\x00", b.Type, len(b.Content))
	data := append([]byte(header), b.Content...)

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	w.Close()
	return buf.Bytes(), nil
}

func (b *BaseObject) Hash() string {
	header := fmt.Sprintf("%s %d\x00", b.Type, len(b.Content))
	data := append([]byte(header), b.Content...)

	hash := sha1.Sum(data)
	return hex.EncodeToString(hash[:])
}
