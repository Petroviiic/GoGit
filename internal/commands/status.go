package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunStatus(repo *core.Repository) error {
	// On branch main
	// Your branch is ahead of 'origin/main' by 3 commits.

	// Changes not staged for commit:
	// (use "git add <file>..." to update what will be committed)
	// (use "git restore <file>..." to discard changes in working directory)
	// 		modified:   cmd/main.go

	// no changes added to commit (use "git add" and/or "git commit -a")

	currentBranch := repo.GetCurrentBranch()
	fmt.Printf("On branch %s", currentBranch)

	workingDirFiles := make(map[string]string)
	if err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		parts := strings.Split(filepath.ToSlash(path), "/")
		if parts[0] == ".git" || parts[0] == ".gogit" {
			return nil
		}

		fullPath := filepath.Join(repo.WorkTree, path)

		info, err := os.Stat(fullPath)
		if info.IsDir() {
			return err
		}

		content, err := os.ReadFile(fullPath)

		blob := core.NewBlob(content)
		workingDirFiles[path] = blob.Hash()

		return err
	}); err != nil {
		return err
	}
	//fmt.Println(workingDirFiles, len(workingDirFiles))

	latestCommit := repo.GetBranchCommit(currentBranch)
	lastCommitFiles := make(map[string]string)

	obj, err := repo.LoadObject(latestCommit)

	commit := obj.(*core.Commit)

	core.GetIndexFromTreeHash(commit.TreeHash, repo, ".", &lastCommitFiles)

	index, err := repo.LoadIndex()

	if err != nil {
		return err
	}

	_ = index
	_ = latestCommit

	return nil
}
