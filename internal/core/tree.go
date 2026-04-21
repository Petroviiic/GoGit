package core

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type TreeEntry struct {
	Mode string // npr. "100644" za običan fajl, "040000" za folder
	Name string // npr. "main.go"
	Hash string // SHA-1 hash tog objekta
}

type Tree struct {
	BaseObject
	Entries []TreeEntry
}

func NewTree(entries []TreeEntry) (*Tree, error) {
	content, err := json.Marshal(entries)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal tree entries: %w", err)
	}

	return &Tree{
		BaseObject: BaseObject{
			Type:    "tree",
			Content: content,
		},
		Entries: entries,
	}, nil
}

type HierarchyNode struct {
	Children map[string]*HierarchyNode // foldername -> node
	Hash     string                    // if file = hash else empty
	IsFile   bool
}

func NewHierarchyNode() *HierarchyNode {
	return &HierarchyNode{
		Children: make(map[string]*HierarchyNode),
	}
}

func CreateFolderHierarchy(index map[string]string) *HierarchyNode {
	root := NewHierarchyNode()

	for path, hash := range index {
		parts := strings.Split(filepath.ToSlash(path), "/")

		currentNode := root

		for i, part := range parts {
			isFile := i == len(parts)-1

			if _, ok := currentNode.Children[part]; !ok {
				node := NewHierarchyNode()
				node.IsFile = isFile

				if isFile {
					node.Hash = hash
				}

				currentNode.Children[part] = node

			}
			currentNode = currentNode.Children[part]
		}
	}

	return root
}

func CreateTreeStructure(folderRoot *HierarchyNode, repo *Repository) (string, error) {
	var dfs func(string, *HierarchyNode) (*TreeEntry, error)
	dfs = func(path string, folderRoot *HierarchyNode) (*TreeEntry, error) {
		if folderRoot == nil {
			return nil, nil
		}
		if folderRoot.IsFile {
			return &TreeEntry{
				Mode: "100644",
				Name: path,
				Hash: folderRoot.Hash,
			}, nil
		}

		var entries []TreeEntry
		for path, node := range folderRoot.Children {
			entry, err := dfs(path, node)
			if err != nil {
				return nil, err
			}
			entries = append(entries, *entry)
		}
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Name < entries[j].Name
		})

		subTree, err := NewTree(entries)
		if err != nil {
			return nil, err
		}

		hashedTree, err := subTree.StoreObject(repo)
		if err != nil {
			return nil, err
		}

		return &TreeEntry{
			Mode: "040000",
			Name: path,
			Hash: hashedTree,
		}, nil
	}

	finalRootEntry, err := dfs(".", folderRoot)
	if err != nil {
		return "", err
	}

	return finalRootEntry.Hash, nil
}
