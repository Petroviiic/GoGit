package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunCommit(repo *core.Repository, message, author string) error {
	fmt.Println(message, author)
	index, err := repo.LoadIndex()

	if err != nil {
		return err
	}

	hierarchyRoot := core.CreateFolderHierarchy(index)

	treeHash, err := core.CreateTreeStructure(hierarchyRoot, repo)
	if err != nil {
		return err
	}

	fmt.Println(treeHash)
	return nil
}
