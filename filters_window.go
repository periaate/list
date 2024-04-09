//go:build windows
// +build windows

package list

import (
	"io/fs"
	"syscall"
)

func addCreationT(fi *Finfo, d fs.DirEntry) bool {
	fileinfo, err := d.Info()
	if err != nil || fileinfo == nil {
		return false
	}
	winFileInfo := fileinfo.Sys().(*syscall.Win32FileAttributeData)

	fi.Vany = winFileInfo.CreationTime.Nanoseconds()

	return true
}
