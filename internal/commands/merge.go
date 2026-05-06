package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunMerge(repo *core.Repository, theirsBranch string) error {
	oursBranch := repo.GetCurrentBranch()
	oursCommit := repo.GetBranchCommit(oursBranch)
	if oursCommit == "" {
		return fmt.Errorf("merge: no commits on branch %s", oursBranch)
	}

	theirsCommit := repo.GetBranchCommit(theirsBranch)
	if theirsCommit == "" {
		return fmt.Errorf("merge: %s - not something we can merge", theirsBranch)
	}

	if oursCommit == theirsCommit {
		fmt.Println("Already up to date.")
		return nil
	}

	// baseCommitHash := repo.GetMergeBase(oursCommitHash, theirsCommitHash)

	return nil
}
