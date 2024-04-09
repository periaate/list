package list

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// This file has largely been generated with GPT4.

type TreeNode struct {
	name     string
	children map[string]*TreeNode
}

func NewTreeNode(name string) *TreeNode {
	return &TreeNode{name: name, children: make(map[string]*TreeNode)}
}

func (t *TreeNode) AddPath(path string) {
	parts := strings.Split(path, "/")
	current := t
	for _, part := range parts {
		if _, exists := current.children[part]; !exists {
			current.children[part] = NewTreeNode(part)
		}
		current = current.children[part]
	}
}

func (t *TreeNode) PrintTree(prefix string) {
	keys := make([]string, 0, len(t.children))
	for k := range t.children {
		keys = append(keys, k)
	}
	sort.Strings(keys) // Sort keys to print in order

	for i, key := range keys {
		child := t.children[key]
		if i == len(t.children)-1 {
			fmt.Println(prefix + "└── " + child.name)
			child.PrintTree(prefix + "    ")
		} else {
			fmt.Println(prefix + "├── " + child.name)
			child.PrintTree(prefix + "│   ")
		}
	}
}

func AddFilesToTree(files []*Finfo) *TreeNode {
	if len(files) == 0 {
		return nil
	}

	commonRoot := strings.Split(filepath.ToSlash(files[0].Path), "/")[0]
	root := NewTreeNode(commonRoot)

	for _, file := range files {
		trimmedPath := strings.TrimPrefix(filepath.ToSlash(file.Path), commonRoot+"/")
		root.AddPath(trimmedPath)
	}

	return root
}
