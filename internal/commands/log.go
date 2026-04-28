package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunLog(limit int, repo *core.Repository) error {
	fmt.Println(limit)
	return nil
}
