package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunInit(repo *core.Repository) error {
	err := repo.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	fmt.Printf("Initialized empty GoGit repository in %s\n", repo.GitDir)
	return nil
}
