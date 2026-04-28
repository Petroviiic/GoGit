package core

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func Deserialize(data []byte) (GitObject, error) {
	r := bytes.NewReader(data)
	close, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer close.Close()

	res, err := io.ReadAll(close)
	if err != nil {
		return nil, err
	}

	index := bytes.Index(res, []byte{'\x00'})

	if index == -1 {
		return nil, fmt.Errorf("error parsing, separator not found")
	}

	header := bytes.Split(res[:index], []byte{' '})
	content := res[index+1:]

	switch string(header[0]) {
	case "blob":
		return NewBlob(content), nil
	case "tree":
		return ParseTree(content)
	case "commit":
		return ParseCommit(content)
	default:
		return nil, fmt.Errorf("unknown object type")
	}
}

func ParseTree(content []byte) (*Tree, error) {
	//fmt.Println("tree content: ", string(content))
	entries := []TreeEntry{}

	err := json.Unmarshal(content, &entries)
	if err != nil {
		return nil, err
	}

	return NewTree(entries)
}

func ParseCommit(content []byte) (*Commit, error) {
	//fmt.Println("commit content: ", string(content))

	lines := bytes.Split(content, []byte("\n"))

	commit := &Commit{}

	messageIndex := 0
	for i, line := range lines {

		if string(line) == "" {
			messageIndex = i + 1
			break
		}

		parts := bytes.SplitN(line, []byte(" "), 2)
		if len(parts) < 2 {
			continue
		}
		key, value := string(parts[0]), parts[1]
		switch key {
		case "tree":
			commit.TreeHash = string(value)
		case "parent":
			commit.ParentHashes = append(commit.ParentHashes, string(value))
		case "author", "committer":
			words := strings.Fields(string(value))
			if len(words) >= 2 {
				timestamp, _ := strconv.ParseInt(words[len(words)-2], 10, 64)

				if key == "author" {
					commit.Author = strings.Join(words[:len(words)-2], " ")
					commit.Timestamp = time.Unix(timestamp, 0).UTC()
				} else {
					commit.Committer = strings.Join(words[:len(words)-2], " ")
				}
			}
		}
	}

	commit.Message = string(bytes.Join(lines[messageIndex:], []byte("\n")))
	return commit, nil
}
