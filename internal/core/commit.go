package core

import (
	"fmt"
	"strings"
	"time"
)

type Commit struct {
	BaseObject
	ParentHashes []string
	Author       string
	Committer    string
	Message      string
	TreeHash     string
	Timestamp    time.Time
}

func NewCommit(author, committer, message, treeHash string, parentHashes []string, timeStamp time.Time) *Commit {
	commit := &Commit{
		BaseObject: BaseObject{
			Type: "commit",
		},
		ParentHashes: parentHashes,
		Author:       author,
		Committer:    committer,
		Message:      message,
		TreeHash:     treeHash,
		Timestamp:    timeStamp,
	}

	content := commit.serialize()

	commit.BaseObject.Content = content

	return commit
}

func (c *Commit) serialize() []byte {
	//tree <tree_hash>
	//parent <parent_hash>
	//author <name> <timestamp> <timezone>
	//committer <name> <timestamp> <timezone>
	//
	//<message>

	lines := []string{fmt.Sprintf("tree %s", c.TreeHash)}
	for _, parentHash := range c.ParentHashes {
		lines = append(lines, fmt.Sprintf("parent %s", parentHash))
	}

	lines = append(lines, fmt.Sprintf("author %s %v +0000", c.Author, c.Timestamp.Unix()))
	lines = append(lines, fmt.Sprintf("committer %s %v +0000", c.Committer, c.Timestamp.Unix()))
	lines = append(lines, "")
	lines = append(lines, c.Message)

	return []byte(strings.Join(lines, "\n"))
}
