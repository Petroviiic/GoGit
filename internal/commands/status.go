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

func RunStatus(repo *core.Repository) error {
	// On branch main
	// Your branch is ahead of 'origin/main' by 3 commits.

	// Changes not staged for commit:
	// (use "git add <file>..." to update what will be committed)
	// (use "git restore <file>..." to discard changes in working directory)
	// 		modified:   cmd/main.go

	// no changes added to commit (use "git add" and/or "git commit -a")

	currentBranch := repo.GetCurrentBranch()
	fmt.Printf("On branch %s\n", currentBranch)

	workingDirFiles := make(map[string]string)
	if err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		parts := strings.Split(filepath.ToSlash(path), "/")
		if slices.Contains(parts, ".git") || slices.Contains(parts, ".gogit") {
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
	if err != nil {
		return err
	}

	commit := obj.(*core.Commit)
	if commit.TreeHash != "" {
		if err := core.GetIndexFromTreeHash(commit.TreeHash, repo, ".", lastCommitFiles); err != nil {
			return err
		}
	}
	// fmt.Println(lastCommitFiles, len(lastCommitFiles))

	index, err := repo.LoadIndex()
	if err != nil {
		return err
	}

	stagedModified := []string{}
	stagedNew := []string{}
	stagedDeleted := []string{}

	for path, hash := range index {
		if lastHash, ok := lastCommitFiles[path]; ok {
			if hash != lastHash {
				stagedModified = append(stagedModified, path)
			}
		} else {
			stagedNew = append(stagedNew, path)
		}
	}

	for path := range lastCommitFiles {
		if _, ok := index[path]; !ok {
			stagedDeleted = append(stagedDeleted, path)
		}
	}

	//fmt.Println(staged)

	//unstaged files
	unstaged := []string{}

	for path, hash := range workingDirFiles {
		lastHash, ok := index[path]
		if ok && hash != lastHash {
			unstaged = append(unstaged, path)
		}
	}

	if len(unstaged) > 0 {
		fmt.Printf("\nChanges not staged for commit:\n\t(use 'git add <file>...' to update what will be committed)\n")
		for _, file := range unstaged {
			fmt.Printf("\t\t%s\n", file)
		}
	}

	//untracked files - Files that are in your directory but have never been added to GoGit's version control.
	untracked := []string{}

	for path := range workingDirFiles {
		_, ok1 := lastCommitFiles[path]
		_, ok2 := index[path]
		if !ok1 && !ok2 {
			untracked = append(untracked, path)
		}
	}
	if len(untracked) > 0 {
		fmt.Printf("\nUntracked files:\n\t(use 'go run ./... add <file>...' to include in what will be committed)\n")
		for _, file := range untracked {
			fmt.Printf("\t\t%s\n", file)
		}
	}

	//deleted
	deleted := []string{}
	for path := range index {
		_, ok := workingDirFiles[path]
		if !ok {
			deleted = append(deleted, path)
		}
	}
	if len(deleted) > 0 {
		fmt.Printf("\nDeleted files:\n")
		for _, file := range deleted {
			fmt.Printf("\t\t%s\n", file)
		}
	}

	return nil
}
