//go:build !windows
// +build !windows

package list

import (
	"io/fs"
)

// Implementation or stub of addCreationT for Unix
func addCreationT(_ *Element, _ fs.FileInfo) {}
