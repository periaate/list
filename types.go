package main

import (
	"archive/zip"
	"io/fs"
	"os"
	"path/filepath"
)

type filter func(*finfo, fs.DirEntry) bool
type process func(filenames []*finfo) []*finfo

type contentType uint8
type sortBy uint8
type orderTo uint8

const (
	other contentType = (1 << iota) - 1
	image
	video
	audio

	byNone sortBy = iota
	byMod
	bySize
	byName

	toDesc orderTo = iota
	toAsc
)

func strToSortBy(s string) sortBy {
	switch s {
	case "date":
		return byMod
	case "mod":
		return byMod
	case "size":
		return bySize
	case "name":
		return byName
	case "none":
		fallthrough
	default:
		return byNone
	}
}

type result struct{ files []*finfo }

type ZipEntry struct{ *zip.File }

func (z ZipEntry) Name() string               { return z.File.Name }
func (z ZipEntry) IsDir() bool                { return z.File.FileInfo().IsDir() }
func (z ZipEntry) Type() fs.FileMode          { return z.File.FileInfo().Mode() }
func (z ZipEntry) Info() (fs.FileInfo, error) { return z.File.FileInfo(), nil }

// check that zipentry is os.DirEntry
var Z os.DirEntry = ZipEntry{}

type finfo struct {
	name string
	path string // includes name, relative path to cwd
	vany int64  // any numeric value, used for sorting
}

func getContentType(filename string) contentType {
	ext := filepath.Ext(filename)
	if t, ok := contentTypes[ext]; ok {
		return t
	}
	return other
}

func stringToContentType(s string) contentType {
	switch s {
	case "image":
		return image
	case "video":
		return video
	case "audio":
		return audio
	default:
		return other
	}
}

var contentTypes = map[string]contentType{
	// image
	".jpg":  image,
	".jpeg": image,
	".png":  image,
	".apng": image,
	".gif":  image,
	".bmp":  image,
	".webp": image,
	".avif": image,
	".jxl":  image,
	".tiff": image,

	// video
	".mp4":  video,
	".m4v":  video,
	".webm": video,
	".mkv":  video,
	".avi":  video,
	".mov":  video,
	".mpg":  video,
	".mpeg": video,

	// audio
	".m4a":  audio,
	".opus": audio,
	".ogg":  audio,
	".mp3":  audio,
	".flac": audio,
	".wav":  audio,
	".aac":  audio,
}
