package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunAdd(args []string, repo *core.Repository) error {
	for _, path := range args {
		return addPath(path, repo)
	}
	return nil
}

func addPath(path string, repo *core.Repository) error {
	index, err := repo.LoadIndex()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(repo.WorkTree, path)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Path %s does not exist", fullPath)
		}
		return err
	}

	if fileInfo.IsDir() {
		err = addDirectory(path, repo, index)
	} else {
		err = addFile(path, repo, index)
	}

	if err != nil {
		return err
	}

	err = repo.SaveIndex(index)
	if err != nil {
		return err
	}

	return nil
}

func addDirectory(path string, repo *core.Repository, index map[string]string) error {
	var files []string
	_ = filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return err
	})

	addedCount := 0
	for _, file := range files {
		if slices.Contains(strings.Split(filepath.ToSlash(file), "/"), ".gogit") {
			continue
		}
		if err := addFile(file, repo, index); err != nil {
			return err
		}
		addedCount++
	}

	fmt.Printf("Added %d files from directory %s", addedCount, path)
	return nil
}

func addFile(path string, repo *core.Repository, index map[string]string) error {
	//reading the file
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	//creating new blob
	blob := core.NewBlob(content)

	//storing under hash[:2]/hash[2:]
	blob_hash, err := blob.StoreObject(repo)
	if err != nil {
		return err
	}

	//updating index
	// relPath, err := filepath.Rel(repo.WorkTree, path)
	// if err != nil {
	// 	return err
	// }
	index[filepath.ToSlash(path)] = blob_hash

	return nil
}
