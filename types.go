package list

import (
	"archive/zip"
	"io/fs"
	"os"
	"path/filepath"
)

type Filter func(*Finfo, fs.DirEntry) bool
type Process func(filenames []*Finfo) []*Finfo

type SortBy uint8
type OrderTo uint8

const (
	Other   = "other"
	Image   = "image"
	Video   = "video"
	Audio   = "audio"
	Archive = "archive"
	ZipLike = "zip"

	ByNone SortBy = iota
	ByMod
	BySize
	ByName

	ToDesc OrderTo = iota
	ToAsc
)

func StrToSortBy(s string) SortBy {
	switch s {
	case "date":
		return ByMod
	case "mod":
		return ByMod
	case "size":
		return BySize
	case "name":
		return ByName
	case "none":
		fallthrough
	default:
		return ByNone
	}
}

type Result struct{ Files []*Finfo }

type ZipEntry struct{ *zip.File }

func (z ZipEntry) Name() string               { return z.File.Name }
func (z ZipEntry) IsDir() bool                { return z.File.FileInfo().IsDir() }
func (z ZipEntry) Type() fs.FileMode          { return z.File.FileInfo().Mode() }
func (z ZipEntry) Info() (fs.FileInfo, error) { return z.File.FileInfo(), nil }

// check that zipentry is os.DirEntry
var _ os.DirEntry = ZipEntry{}

type Finfo struct {
	name string
	path string // includes name, relative path to cwd
	vany int64  // any numeric value, used for sorting
}

func GetContentTypes(filename string) (res ArrSet[string]) {
	ext := filepath.Ext(filename)
	for k, v := range CntType {
		if v.Contains(ext) {
			res = append(res, k)
		}
	}
	return
}

type ArrSet[T comparable] []T

func (a ArrSet[T]) Contains(ext T) bool {
	for _, v := range a {
		if v == ext {
			return true
		}
	}
	return false
}

var CntType = map[string]ArrSet[string]{
	Image:   {".jpg", ".jpeg", ".png", ".apng", ".gif", ".bmp", ".webp", ".avif", ".jxl", ".tiff"},
	Video:   {".mp4", ".m4v", ".webm", ".mkv", ".avi", ".mov", ".mpg", ".mpeg"},
	Audio:   {".m4a", ".opus", ".ogg", ".mp3", ".flac", ".wav", ".aac"},
	Archive: {".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".lz4", ".zst", ".lzma", ".lzip", ".lz", ".cbz"},
	ZipLike: {".zip", ".cbz", ".cbr"},
}
