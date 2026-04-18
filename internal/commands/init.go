package commands

import (
	"fmt"
	"path/filepath"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunInit(args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	repo := core.NewRepository(absPath)
	err = repo.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}

	fmt.Printf("Initialized empty GoGit repository in %s\n", repo.GitDir)
	return nil
}
