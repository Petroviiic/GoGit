package core

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

type GitObject interface {
	GetType() string
	GetContent() []byte
	Serialize() ([]byte, error)
	Hash() string
	StoreObject(*Repository) (string, error)
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

func (b *BaseObject) StoreObject(repo *Repository) (string, error) {
	hash := b.Hash()

	objDir := filepath.Join(repo.ObjectsDir, hash[:2])
	objFile := filepath.Join(objDir, hash[2:])

	content, err := b.Serialize()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(objDir, 0755)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(objFile, content, 0644); err != nil {
		return "", err
	}
	return hash, nil
}
