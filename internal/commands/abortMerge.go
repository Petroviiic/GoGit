package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func AbortMerge(repo *core.Repository) error {
	fmt.Println("aborting merge")

	//check if merge head exists
	if _, exists := repo.MergeHEADExists(); !exists {
		return fmt.Errorf("fatal: There is no merge to abort") //promijeni u mainu da ispise fatal a ne error :)
	}

	currentBranch, isDetached := repo.GetCurrentBranch()

	if isDetached {
		return fmt.Errorf("branch detached. something went terribly wrong!!!")
	}

	// RestoreWorkingDirectoryFiles()
	lastOursCommitHash := repo.GetBranchCommit(currentBranch)
	obj, err := repo.LoadObject(lastOursCommitHash)
	if err != nil {
		return err
	}

	lastOursCommit := obj.(*core.Commit)
	lastOursIndex := map[string]core.IndexEntry{}
	if err := core.GetFilesFromTreeHash(lastOursCommit.TreeHash, repo, "", lastOursIndex); err != nil {
		return err
	}

	index, err := repo.LoadIndex()

	if err != nil {
		return err
	}

	//restore working dir and index to head.lastcommit
	if err := RestoreWorkingDirectoryFiles(lastOursIndex, index, "", repo); err != nil {
		return err
	}

	if err := repo.SaveIndex(lastOursIndex); err != nil {
		return err
	}

	//remove merge head
	if err := repo.DeleteMergeHEAD(); err != nil {
		return err
	}

	fmt.Println("merge sucessfully aborted")
	return nil
}
