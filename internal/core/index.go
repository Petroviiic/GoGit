package core

import (
	"encoding/json"
	"os"
)

type IndexEntry struct {
	Hash  string `json:"hash"`
	MTime int64  `json:"mtime"`
}

func (r *Repository) LoadIndex() (map[string]IndexEntry, error) {
	indexRaw, err := os.ReadFile(r.IndexPath)
	if err != nil {
		return nil, err
	}
	//reading index
	var index map[string]IndexEntry
	err = json.Unmarshal(indexRaw, &index)
	if err != nil {
		return nil, err
	}

	return index, err
}

func (r *Repository) SaveIndex(index map[string]IndexEntry) error {
	newIndex, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(r.IndexPath, newIndex, 0644)
	if err != nil {
		return err
	}

	return nil
}
