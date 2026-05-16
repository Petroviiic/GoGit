package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func AbortMerge(repo *core.Repository) error {
	fmt.Println("aborting merge")
	return nil
}
