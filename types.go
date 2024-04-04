package list

import (
	"archive/zip"
	"io/fs"
	"os"
	"path/filepath"
)

type filter func(*Finfo, fs.DirEntry) bool
type process func(filenames []*Finfo) []*Finfo

type sortBy uint8
type orderTo uint8

const (
	other   = "other"
	image   = "image"
	video   = "video"
	audio   = "audio"
	archive = "archive"
	zipLike = "zip"

	byNone sortBy = iota
	byMod
	bySize
	byName

	toDesc orderTo = iota
	toAsc
)

func StrToSortBy(s string) sortBy {
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

type Result struct{ Files []*Finfo }

type ZipEntry struct{ *zip.File }

func (z ZipEntry) Name() string               { return z.File.Name }
func (z ZipEntry) IsDir() bool                { return z.File.FileInfo().IsDir() }
func (z ZipEntry) Type() fs.FileMode          { return z.File.FileInfo().Mode() }
func (z ZipEntry) Info() (fs.FileInfo, error) { return z.File.FileInfo(), nil }

// check that zipentry is os.DirEntry
var Z os.DirEntry = ZipEntry{}

type Finfo struct {
	name string
	path string // includes name, relative path to cwd
	vany int64  // any numeric value, used for sorting
}

func GetContentTypes(filename string) (res arrSet[string]) {
	ext := filepath.Ext(filename)
	for k, v := range cntType {
		if v.contains(ext) {
			res = append(res, k)
		}
	}
	return
}

type arrSet[T comparable] []T

func (a arrSet[T]) contains(ext T) bool {
	for _, v := range a {
		if v == ext {
			return true
		}
	}
	return false
}

var cntType = map[string]arrSet[string]{
	image:   {".jpg", ".jpeg", ".png", ".apng", ".gif", ".bmp", ".webp", ".avif", ".jxl", ".tiff"},
	video:   {".mp4", ".m4v", ".webm", ".mkv", ".avi", ".mov", ".mpg", ".mpeg"},
	audio:   {".m4a", ".opus", ".ogg", ".mp3", ".flac", ".wav", ".aac"},
	archive: {".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".lz4", ".zst", ".lzma", ".lzip", ".lz", ".cbz"},
	zipLike: {".zip", ".cbz", ".cbr"},
}
