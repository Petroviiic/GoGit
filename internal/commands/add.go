package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunAdd(args []string, repo *core.Repository) error {
	for _, path := range args {
		return addPath(path, repo)
	}
	return nil
}

func addPath(path string, repo *core.Repository) error {
	fullPath := filepath.Join(repo.WorkTree, path)

	fmt.Println(fullPath, path, repo.GitDir, repo.WorkTree)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Path %s does not exist", fullPath)
		}
		return err
	}

	if fileInfo.IsDir() {
		return addDirectory(path, repo)
	} else {
		return addFile(path, repo)
	}
}

func addDirectory(path string, repo *core.Repository) error {

	return nil
}

func addFile(path string, repo *core.Repository) error {
	//reading the file
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	//creating new blob
	blob := core.NewBlob(content)

	//storing under hash[:2]/hash[2:]
	blob_hash, err := StoreObject(&blob.BaseObject, repo)
	if err != nil {
		return err
	}

	//updating index
	indexRaw, err := os.ReadFile(repo.IndexPath)
	if err != nil {
		return err
	}
	//reading index
	var index map[string]string
	err = json.Unmarshal(indexRaw, &index)
	if err != nil {
		return err
	}
	//updating index
	index[path] = blob_hash[2:]

	newIndex, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	//saving changes
	err = os.WriteFile(repo.IndexPath, newIndex, 0644)
	if err != nil {
		return err
	}
	return nil
}

func StoreObject(o *core.BaseObject, repo *core.Repository) (string, error) {
	hash := o.Hash()

	objDir := filepath.Join(repo.ObjectsDir, hash[:2])
	objFile := filepath.Join(objDir, hash[2:])

	content, err := o.Serialize()
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
