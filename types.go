package list

import (
	"archive/zip"
	"io/fs"
	"os"
)

type Filter func(*Finfo, fs.DirEntry) bool
type Process func(filenames []*Finfo) []*Finfo

type SortBy uint8
type OrderTo uint8

func AsMask(sar []string) uint32 {
	var mask uint32
	for _, v := range sar {
		mask |= StrToMask(v)
	}
	return mask
}

func StrToMask(str string) uint32 {
	switch str {
	case Image:
		return MaskImage
	case Video:
		return MaskVideo
	case Audio:
		return MaskAudio
	case Archive:
		return MaskArchive
	case ZipLike:
		return MaskZipLike
	default:
		return 0
	}
}

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
	ByCreation
	ByName

	ToDesc OrderTo = iota
	ToAsc

	_ uint32 = 1 << iota
	MaskImage
	MaskVideo
	MaskAudio
	MaskArchive
	MaskZipLike = 1<<iota + MaskArchive
)

func StrToSortBy(s string) SortBy {
	switch s {
	case "date":
		return ByMod
	case "mod":
		return ByMod
	case "creation":
		return ByCreation
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

func (r Result) Sar() []string {
	res := make([]string, 0, len(r.Files))
	for _, v := range r.Files {
		res = append(res, v.Path)
	}
	return res
}

type ZipEntry struct{ *zip.File }

func (z ZipEntry) Name() string               { return z.File.Name }
func (z ZipEntry) IsDir() bool                { return z.File.FileInfo().IsDir() }
func (z ZipEntry) Type() fs.FileMode          { return z.File.FileInfo().Mode() }
func (z ZipEntry) Info() (fs.FileInfo, error) { return z.File.FileInfo(), nil }

// check that zipentry is os.DirEntry
var _ os.DirEntry = ZipEntry{}

type Finfo struct {
	Name      string
	Path      string // includes name, relative path to cwd
	Vany      int64  // any numeric value, used for sorting
	Mask      uint32 // file kind, bitmask, see Mask* constants
	IsDir     bool
	IsArchive bool
}

var CntMasks = map[uint32][]string{
	MaskImage:   {".jpg", ".jpeg", ".png", ".apng", ".gif", ".bmp", ".webp", ".avif", ".jxl", ".tiff"},
	MaskVideo:   {".mp4", ".m4v", ".webm", ".mkv", ".avi", ".mov", ".mpg", ".mpeg"},
	MaskAudio:   {".m4a", ".opus", ".ogg", ".mp3", ".flac", ".wav", ".aac"},
	MaskArchive: {".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".lz4", ".zst", ".lzma", ".lzip", ".lz", ".cbz"},
	MaskZipLike: {".zip", ".cbz", ".cbr"},
}

func RegisterMasks(mask uint32, keys ...string) {
	for _, k := range keys {
		CntMap[k] |= mask
	}
}

var CntMap = map[string]uint32{}

func init() {
	for k, v := range CntMasks {
		RegisterMasks(k, v...)
	}
}
