package list

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const (
	Other   = "other"
	Image   = "image"
	Video   = "video"
	Audio   = "audio"
	Archive = "archive"
	ZipLike = "zip"

	_ uint32 = 1 << iota
	MaskImage
	MaskVideo
	MaskAudio
	MaskArchive
	MaskZipLike = 1<<iota + MaskArchive
)

var CntMasks = map[uint32][]string{
	MaskImage:   {".jpg", ".jpeg", ".png", ".apng", ".gif", ".bmp", ".webp", ".avif", ".jxl", ".tiff"},
	MaskVideo:   {".mp4", ".m4v", ".webm", ".mkv", ".avi", ".mov", ".mpg", ".mpeg"},
	MaskAudio:   {".m4a", ".opus", ".ogg", ".mp3", ".flac", ".wav", ".aac"},
	MaskArchive: {".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".lz4", ".zst", ".lzma", ".lzip", ".lz", ".cbz"},
	MaskZipLike: {".zip", ".cbz", ".cbr"},
}
var CntMap = map[string]uint32{}

// Traverse traverses directories non-recursively and breadth first.
func Traverse(wfn fs.WalkDirFunc, opts *Options) {
	dirs := opts.Args
	var depth int
	for len(dirs) != 0 {
		if depth > opts.ToDepth {
			return
		}
		var nd []string
		for _, d := range dirs {
			ext := filepath.Ext(d)
			slog.Debug("traversing", "dir", d, "depth", depth, "ext", ext, "isarchive", CntMap[ext]&MaskZipLike != 0)

			if opts.Archive && CntMap[ext]&MaskZipLike != 0 {
				TraverseZip(d, depth, wfn, opts)
				continue
			}

			entries, err := os.ReadDir(d)
			if err != nil {
				slog.Debug("found a non-directory argument", "argument", d)
				continue
			}
			for _, entry := range entries {
				path := filepath.Join(d, entry.Name())
				if entry.IsDir() {
					nd = append(nd, path)
				}

				if opts.Archive && filepath.Ext(path) == ".zip" {
					nd = append(nd, path)
					continue
				}

				if depth < opts.FromDepth {
					continue
				}

				err := wfn(path, entry, nil)
				if err != nil {
					continue
				}
			}
		}

		dirs = nd
		depth++
	}
}

func TraverseZip(path string, depth int, wfn fs.WalkDirFunc, opts *Options) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Close()

	for _, f := range r.File {
		fn := filepath.ToSlash(f.Name)

		fdepth := depth + strings.Count(fn, "/")
		if fdepth < opts.FromDepth || fdepth > opts.ToDepth {
			continue
		}
		fpath := filepath.Join(path, fn)

		entry := ZipEntry{f}

		if depth <= opts.FromDepth {
			continue
		}

		err := wfn(fpath, entry, fmt.Errorf("zl"))
		if err != nil {
			continue
		}
	}
}

func BuildWalkDirFn(fns []Filter, res *Result) func(string, fs.DirEntry, error) error {
	return func(path string, d fs.DirEntry, err error) error {
		var b bool
		if err != nil {
			r := err.Error()
			if r == "zl" {
				err = nil
				b = true
			}
		}

		if d == nil || err != nil {
			return nil
		}

		fi := &Finfo{Name: d.Name(), Path: path, IsDir: d.IsDir()}

		ext := filepath.Ext(path)
		fi.Mask |= CntMap[ext]
		if b {
			fi.Mask |= MaskZipLike
			fi.IsArchive = true
		}

		for _, fn := range fns {
			res := fn(fi, d)
			if !res {
				return nil
			}
		}
		res.Files = append(res.Files, fi)
		return nil
	}
}

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

func RegisterMasks(mask uint32, keys ...string) {
	for _, k := range keys {
		CntMap[k] |= mask
	}
}

func init() {
	for k, v := range CntMasks {
		RegisterMasks(k, v...)
	}
}
