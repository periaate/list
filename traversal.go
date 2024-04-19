package list

import (
	"archive/zip"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
)

type Element struct {
	Name      string
	Path      string // includes name, relative path to cwd
	Vany      int64  // any numeric value, used for sorting
	Mask      uint32 // file kind, bitmask, see Mask* constants
	IsDir     bool
	IsArchive bool // is a readable archive; ziplike
}

type ResultFn func(*Element)

type Parser func(string, fs.FileInfo) *Element

func InitFileParser(opts *Options) Parser {
	return func(path string, info fs.FileInfo) *Element {
		if !opts.NoHide {
			if _, ok := Hide[info.Name()]; ok || info.Name()[0] == '.' {
				return nil
			}
		}
		path = filepath.ToSlash(path)
		path = filepath.Join(path, info.Name())
		fi := &Element{
			Name:  info.Name(),
			Path:  path,
			IsDir: info.IsDir(),
		}

		fi.Mask |= CntMap[filepath.Ext(fi.Name)]

		if fi.Mask&MaskZipLike != 0 {
			fi.IsArchive = true
		}

		switch StrToSortBy(opts.Sort) {
		case ByMod:
			addModT(fi, info)
		case BySize:
			addSize(fi, info)
		case ByCreation:
			addCreationT(fi, info)
		}

		return fi
	}
}
func addModT(fi *Element, info fs.FileInfo) { fi.Vany = info.ModTime().Unix() }
func addSize(fi *Element, info fs.FileInfo) { fi.Vany = info.Size() }

type Yield func(paths []string) (el []*Element, ok bool)

func ResolveHome(path string) string {
	if len(path) == 0 {
		return path
	}
	if path[0] == '~' {
		dirname, err := os.UserHomeDir()
		if err != nil {
			slog.Error("error resolving home dir", "error", err)
			return path
		}
		path = filepath.Join(dirname, path[1:])
	}
	return path
}

func Traverse(opts *Options, yfn Yield, rfn ResultFn) {
	var depth int
	dirPaths := opts.Args
	for i, path := range dirPaths {
		dirPaths[i] = ResolveHome(path)
	}
	slog.Debug("traversing", "dirs", dirPaths)

	for els, ok := yfn(dirPaths); ok; els, ok = yfn(dirPaths) {
		dirPaths = make([]string, 0)
		if depth > opts.ToDepth {
			slog.Debug("reached max depth", "depth", depth, "todepth", opts.ToDepth)
			return
		}

		for _, el := range els {
			if el.IsDir {
				dirPaths = append(dirPaths, el.Path)
				if opts.OnlyFiles {
					continue
				}
				if opts.DirOnly {
					rfn(el)
					continue
				}
			}

			if depth < opts.FromDepth {
				slog.Debug("skipping", "element", el.Path, "depth", depth)
				continue
			}

			rfn(el)
		}

		depth++
	}
}

func GetYieldFs(opts *Options) Yield {
	parser := InitFileParser(opts)
	return func(paths []string) (els []*Element, ok bool) {
		slog.Debug("yielding", "paths", len(paths))
		if len(paths) == 0 {
			return nil, false
		}
		for _, path := range paths {
			var err error
			var finfos []fs.FileInfo
			switch {
			case opts.Archive && IsZipLike(path):
				finfos, err = TraverseZip(path)
			default:
				finfos, err = TraverseDir(path)
			}
			if err != nil {
				slog.Debug("error during traversal", "err", err)
				continue
			}

			for _, finfo := range finfos {
				el := parser(path, finfo)
				if el != nil {
					els = append(els, el)
				}
			}

		}
		if len(els) == 0 {
			slog.Debug("no elements found")
			return nil, false
		}
		ok = true
		slog.Debug("yielding", "elements", len(els))
		return
	}
}

func TraverseDir(path string) (files []fs.FileInfo, err error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			slog.Debug("error reading file info", "file", entry.Name(), "error", err)
			continue
		}
		files = append(files, info)
	}
	return
}

func TraverseZip(path string) (files []fs.FileInfo, err error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		info := f.FileInfo()
		if info.IsDir() {
			continue
		}

		files = append(files, info)
	}

	return
}

func IsZipLike(path string) bool { return CntMap[filepath.Ext(path)]&MaskZipLike != 0 }

func GetRfn(f Filter, res *Result) ResultFn {
	return func(el *Element) {
		if !f(el) {
			return
		}
		res.Files = append(res.Files, el)
	}
}
