package list

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// the interval of flushing the buffer
const bufLength = 500

func PrintWithBuf(els []*Element, opts *Options) {
	if opts.Count {
		fmt.Println(len(els))
		return
	}

	if len(els) == 0 {
		return
	}
	if opts.Quiet {
		slog.Debug("quiet flag is set, returning from print function")
		return
	}

	if opts.Tree {
		ftree := AddFilesToTree(els)
		ftree.PrintTree("")
		return
	}

	// I am unsure of how large this buffer should be. Testing or profiling might be necessary to
	// find what is reasonable. The default buffer size was flushing automatically before being told to.
	// This might be okay in itself, and we might not need to manually set a buffer ta all (or flush).

	w := bufio.NewWriterSize(os.Stdout, 4096*bufLength)

	for i, file := range els {
		fp := filepath.ToSlash(file.Path)
		if opts.Absolute {
			fp, _ = filepath.Abs(file.Path)
			fp = filepath.ToSlash(fp)
		}
		res := fp + "\n"

		w.WriteString(res)
		if i%bufLength == 0 {
			w.Flush()
		}
	}

	w.Flush()
}

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

func AddFilesToTree(files []*Element) *TreeNode {
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
