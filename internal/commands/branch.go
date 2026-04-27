package commands

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunBranch(branch string, shouldDelete, listOnly bool, repo *core.Repository) error {
	fmt.Printf("branch: %s, shouldDelete %v, listOnly: %v \n", branch, shouldDelete, listOnly)

	branches := []string{}
	if err := filepath.Walk(repo.RefsDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		_, file := filepath.Split(path)
		branches = append(branches, file)
		return nil
	}); err != nil {
		return err
	}

	current := repo.GetCurrentBranch()

	//<no_param>   	-lists all the branches
	if listOnly {
		fmt.Printf("* %s\n", current)

		for _, branch := range branches {
			if branch == current {
				continue
			}
			fmt.Printf("%s\n", branch)
		}
		return nil
	}

	//-d <branch_name>   	-deletes the branch
	if shouldDelete {
		if !slices.Contains(branches, branch) {
			return fmt.Errorf("error: branch %s not found", branch)
		}

		if err := os.Remove(filepath.Join(repo.RefsDir, branch)); err != nil {
			return err
		}

		return nil
	}

	//<branch_name>			-creates a new branch
	if slices.Contains(branches, branch) {
		return fmt.Errorf("a branch named %s already exists", branch)
	}

	lastCommit := repo.GetBranchCommit(current)
	if lastCommit == "" {
		return fmt.Errorf("no commits yet, cannot create a branch")
	}

	if err := repo.SetBranchCommit(branch, lastCommit); err != nil { //if its brand new, it should point to its parent last commit
		return err
	}
	fmt.Printf("created new branch %s", branch)
	return nil
}
