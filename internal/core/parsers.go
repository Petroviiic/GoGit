package core

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
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
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown object type")
	}
}

func ParseTree(content []byte) (*Tree, error) {
	fmt.Println("tree content: ", string(content))
	entries := []TreeEntry{}

	err := json.Unmarshal(content, &entries)
	if err != nil {
		return nil, err
	}

	return NewTree(entries)
}
