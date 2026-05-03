package commands

import (
	"fmt"
	"time"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunCommit(repo *core.Repository, message, author string) error {
	index, err := repo.LoadIndex()

	if err != nil {
		return err
	}

	currentBranch := repo.GetCurrentBranch()
	if len(index) == 0 {
		fmt.Printf("nothing to commit. staging area empty")
		return nil
	}
	branchCommit := repo.GetBranchCommit(currentBranch)

	parentHashes := []string{branchCommit}

	hierarchyRoot := core.CreateFolderHierarchy(index)

	treeHash, err := core.CreateTreeStructure(hierarchyRoot, repo)

	if err != nil {
		return err
	}

	if len(parentHashes) > 0 {
		for _, hash := range parentHashes {
			if hash == "" {
				fmt.Println("hash empty, skipping")
				continue
			}
			c, err := repo.LoadObject(hash)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if c.(*core.Commit).TreeHash == treeHash {
				fmt.Printf("Your branch is up to date with '%s'.", currentBranch)
				return nil

			}
		}
	}
	commit := core.NewCommit(
		author,
		author,
		message,
		treeHash,
		parentHashes,
		time.Now().UTC(),
	)

	commitHash, err := commit.StoreObject(repo)
	if err != nil {
		return err
	}

	err = repo.SetBranchCommit(currentBranch, commitHash)
	if err != nil {
		return err
	}

	// err = repo.SaveIndex(make(map[string]string))
	// if err != nil {
	// 	return err
	// }

	fmt.Printf("Created commit %s on branch %s", commitHash, currentBranch)
	return nil
}
