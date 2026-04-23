package commands

import (
	"fmt"
	"time"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunCommit(repo *core.Repository, message, author string) error {
	fmt.Println(message, author)
	// index, err := repo.LoadIndex()

	// if err != nil {
	// 	return err
	// }

	// hierarchyRoot := core.CreateFolderHierarchy(index)

	// treeHash, err := core.CreateTreeStructure(hierarchyRoot, repo)
	// if err != nil {
	// 	return err
	// }
	//fmt.Println(treeHash)

	currentBranch := repo.GetCurrentBranch()
	branchCommit := repo.GetBranchCommit(currentBranch)
	parentHashes := []string{branchCommit}

	commit := core.NewCommit(
		author,
		author,
		message,
		"treeHash",
		parentHashes,
		time.Now().UTC(),
	)
	_ = commit

	// deserialized, _ := core.ParseCommit(commit.Content)

	// fmt.Println(
	// 	"printam",
	// 	deserialized.Author,
	// 	deserialized.Committer,
	// 	deserialized.Message,
	// 	deserialized.ParentHashes,
	// 	deserialized.Timestamp,
	// 	string(deserialized.BaseObject.Content),
	// )
	//fmt.Println(commit.Serialize())

	return nil
}
