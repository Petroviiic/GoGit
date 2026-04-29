package commands

import (
	"fmt"

	"github.com/Petroviiic/GoGit/internal/core"
)

func RunLog(limit int, repo *core.Repository) error {
	currentBranch := repo.GetCurrentBranch()
	latestCommit := repo.GetBranchCommit(currentBranch)

	if latestCommit == "" {
		return fmt.Errorf("no commits on current branch")
	}
	commits := []string{latestCommit}
	visited := make(map[string]bool)
	count := 0

	for len(commits) > 0 && count < limit {
		currentHash := commits[0]

		commits = commits[1:]

		if visited[currentHash] {
			continue
		}

		commitObj, err := repo.LoadObject(currentHash)
		if err != nil {
			fmt.Println(err)
			continue

		}
		commit := commitObj.(*core.Commit)
		fmt.Printf("commit %s\nAuthor: %s\nDate: %v\n\n%s\n", currentHash, commit.Author, commit.Timestamp, commit.Message)

		visited[currentHash] = true
		count++
		for _, hash := range commit.ParentHashes {
			if !visited[hash] {
				commits = append(commits, hash)
			}
		}

	}
	return nil
}
