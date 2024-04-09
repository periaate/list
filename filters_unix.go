//go:build !windows
// +build !windows

package list

import (
	"io/fs"
)

// Implementation or stub of addCreationT for Unix
func addCreationT(fi *Finfo, d fs.DirEntry) bool { return false }
