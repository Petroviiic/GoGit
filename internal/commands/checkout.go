package commands

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunCheckout(branch string, shouldCreate bool, repo *core.Repository) error {
	fmt.Println(branch, shouldCreate)

	lastBranch := repo.GetCurrentBranch()
	lastCommit := repo.GetBranchCommit(lastBranch)

	commitObj, err := repo.LoadObject(lastCommit)
	if err != nil {
		return err
	}

	lastBranchFiles := []string{}

	if err := core.GetFilesFromTreeHash(commitObj.(*core.Commit).TreeHash, repo, "", &lastBranchFiles); err != nil {
		return err
	}

	fmt.Println(lastBranchFiles)

	branchFile := filepath.Join(repo.RefsDir, branch)
	if _, err := os.Stat(branchFile); errors.Is(err, os.ErrNotExist) { //branch doesnt exist
		if shouldCreate {
			if lastCommit == "" {
				return fmt.Errorf("no commits yet, cannot create a branch")
			} else {
				if err := repo.SetBranchCommit(branch, lastCommit); err != nil { //if its brand new, it should point to its parent last commit
					return err
				}
				fmt.Printf("created new branch %s", branch)
			}
		} else {
			return fmt.Errorf("branch %s not found\nuse checkout -b %s to create and switch to a new branch", branch, branch)
		}
	}

	if err := repo.SetCurrentBranch(branch); err != nil {
		return err
	}

	if err := RemoveOldFiles(lastBranchFiles); err != nil {
		return err
	}

	newBranchCommitHash := repo.GetBranchCommit(branch)
	newBranchCommit, err := repo.LoadObject(newBranchCommitHash)

	if err != nil {
		return err
	}

	if err := RestoreWorkingDirectoryFiles(newBranchCommit.(*core.Commit).TreeHash); err != nil {
		return err
	}

	fmt.Println(lastBranch, lastCommit)

	fmt.Printf("switched to branch %s", branch)
	return nil
}

func RemoveOldFiles(any) error {
	return nil
}
func RestoreWorkingDirectoryFiles(any) error {
	return nil
}
