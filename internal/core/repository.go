package core

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Repository struct {
	WorkTree string
	GitDir   string

	ObjectsDir string
	RefsDir    string
	IndexPath  string
	HeadPath   string
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

	repo := &Repository{
		WorkTree:   absPath,
		GitDir:     gitDir,
		ObjectsDir: filepath.Join(gitDir, "objects"),
		RefsDir:    filepath.Join(gitDir, "refs/heads"),
		IndexPath:  filepath.Join(gitDir, "index"),
		HeadPath:   filepath.Join(gitDir, "HEAD"),
	}

	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return repo, fmt.Errorf("repository not found: .gogit")
	}
	return repo, nil
}

func (r *Repository) Init() error {
	if _, err := os.Stat(r.GitDir); !os.IsNotExist(err) {
		return fmt.Errorf("repository in %s already exists", r.GitDir)
	}

	dirs := []string{
		r.GitDir,
		filepath.Join(r.GitDir, "objects"),
		filepath.Join(r.GitDir, "refs"),
		filepath.Join(r.GitDir, "refs/heads"),
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

func (r *Repository) GetCurrentBranch() (string, bool) {
	data, err := os.ReadFile(r.HeadPath)
	if err != nil {
		return "main", false
	}
	content := strings.TrimSpace(string(data))

	if strings.HasPrefix(content, "ref: refs/heads/") {
		return strings.TrimPrefix(content, "ref: refs/heads/"), false
	}

	return content, true //detached head
}

func (r *Repository) SetCurrentBranch(newBranch string, isDetached bool) error {
	ref := ""
	if !isDetached {
		ref = fmt.Sprintf("ref: refs/heads/%s\n", newBranch)
	} else {
		ref = fmt.Sprintf("%s\n", newBranch)
	}

	return os.WriteFile(r.HeadPath, []byte(ref), 0644)
}

func (r *Repository) GetBranchCommit(branch string) string {
	if _, err := os.Stat(r.RefsDir); os.IsNotExist(err) {
		return ""
	}

	branchFile := filepath.Join(r.RefsDir, branch)

	if _, err := os.Stat(branchFile); errors.Is(err, fs.ErrNotExist) {
		return ""
	}

	data, err := os.ReadFile(branchFile)

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

func (repo *Repository) LoadObject(objectHash string) (GitObject, error) {
	if objectHash == "" {
		return nil, fmt.Errorf("couldnt load object; hash empty %s", objectHash)
	}
	objDir := filepath.Join(repo.ObjectsDir, objectHash[:2])

	objFile := filepath.Join(objDir, objectHash[2:])

	data, err := os.ReadFile(objFile)
	if err != nil {
		return nil, err
	}

	return Deserialize(data)
}

func (repo *Repository) SetBranchCommit(branch, hash string) error {
	branchFile := filepath.Join(repo.RefsDir, branch)

	err := os.WriteFile(branchFile, []byte(hash+"\n"), 0644)

	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) GetMergeBase(oursCommitHash, theirsCommitHash string) (string, error) {
	oursCommitHistory := map[string]bool{}

	queue := []string{oursCommitHash}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		oursCommitHistory[current] = true

		obj, err := repo.LoadObject(current)
		if err != nil || obj == nil {
			break
		}
		oursCommit := obj.(*Commit)

		for _, parent := range oursCommit.ParentHashes {
			if !oursCommitHistory[parent] && parent != "" {
				queue = append(queue, parent)
			}
		}
	}

	theirsCommitHistory := map[string]bool{}

	queue = []string{theirsCommitHash}

	result := ""
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		theirsCommitHistory[current] = true

		if oursCommitHistory[current] {
			result = current
			break
		}
		obj, err := repo.LoadObject(current)
		if err != nil || obj == nil {
			break
		}
		theirsCommit := obj.(*Commit)

		for _, parent := range theirsCommit.ParentHashes {
			if !theirsCommitHistory[parent] && parent != "" {
				queue = append(queue, parent)
			}
		}
	}

	if result == "" {
		return "", fmt.Errorf("something went wrong, no merge base found")
	}
	return result, nil
}
