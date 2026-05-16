package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func AbortMerge(repo *core.Repository) error {
	fmt.Println("aborting merge")

	//check if merge head exists
	//remove merge head
	//restore working dir and index to head.lastcommit
	return nil
}
