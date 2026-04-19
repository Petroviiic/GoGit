package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunCommit(repo *core.Repository, message, author string) error {
	fmt.Println(message, author)
	return nil
}
