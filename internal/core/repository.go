package core

import (
	"fmt"
	"os"
	"path/filepath"
)

type Repository struct {
	WorkTree string
	GitDir   string

	ObjectsDir string
	RefsDir    string
	IndexPath  string
}

func NewRepository(args []string) (*Repository, error) {
	path := "."
	// if len(args) > 0 {
	// 	path = args[0]
	// }

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	gitDir := filepath.Join(path, ".gogit")

	return &Repository{
		WorkTree:   absPath,
		GitDir:     gitDir,
		ObjectsDir: filepath.Join(gitDir, "objects"),
		RefsDir:    filepath.Join(gitDir, "refs"),
		IndexPath:  filepath.Join(gitDir, "index"),
	}, nil
}

func (r *Repository) Init() error {
	if _, err := os.Stat(r.GitDir); !os.IsNotExist(err) {
		return fmt.Errorf("repository in %s already exists", r.GitDir)
	}

	dirs := []string{
		r.GitDir,
		filepath.Join(r.GitDir, "objects"),
		filepath.Join(r.GitDir, "refs"),
		filepath.Join(r.GitDir, "refs/dirs"),
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	headPath := filepath.Join(r.GitDir, "HEAD")
	if err := os.WriteFile(headPath, []byte("ref: refs/heads/main\n"), 0644); err != nil {
		return err
	}

	indexPath := filepath.Join(r.GitDir, "index")
	if err := os.WriteFile(indexPath, []byte("{}"), 0644); err != nil {
		return err
	}
	return nil
}
