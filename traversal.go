package list

import (
	"archive/zip"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Result struct{ Files []*Finfo }

func (r Result) Sar() []string {
	res := make([]string, 0, len(r.Files))
	for _, v := range r.Files {
		res = append(res, v.Path)
	}
	return res
}

type Finfo struct {
	Name      string
	Path      string // includes name, relative path to cwd
	Vany      int64  // any numeric value, used for sorting
	Mask      uint32 // file kind, bitmask, see Mask* constants
	IsDir     bool
	IsArchive bool // is a readable archive; ziplike
}

type ResultFilters func(*Finfo)

type FinfoParser func(fs.FileInfo) *Finfo

func InitFileParser(opts *Options) FinfoParser {
	return func(info fs.FileInfo) *Finfo {
		fi := &Finfo{
			Name:  info.Name(),
			Path:  info.Name(),
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
func addModT(fi *Finfo, info fs.FileInfo) { fi.Vany = info.ModTime().Unix() }
func addSize(fi *Finfo, info fs.FileInfo) { fi.Vany = info.Size() }

type Traverser func(*Options, ResultFilters)

func GetTraverser(opts *Options) Traverser {
	switch {
	case opts.ArgMode:
		return TraverseArgs
	case opts.FileMode != "":
		return FileTraverser
	default:
		return TraverseFS
	}
}

func FileTraverser(opts *Options, rfn ResultFilters) {
	for _, arg := range opts.Args {
		b, err := os.ReadFile(arg)
		if err != nil {
			slog.Error("error reading file", "arg", arg, "error", err)
			continue
		}

		var res []string

		switch opts.FileMode {
		case "words", "w":
			res = strings.Fields(string(b))
		case "lines", "l":
			res = strings.Split(string(b), "\n")
		default:
			res = strings.Split(string(b), "\n")
		}

		for _, line := range res {
			rfn(StringParser(line))
		}
	}
}

func TraverseArgs(opts *Options, rfn ResultFilters) {
	for _, arg := range opts.Args {
		rfn(StringParser(arg))
	}
}

func StringParser(s string) *Finfo {
	return &Finfo{
		Name: s,
		Path: s,
		Mask: CntMap[filepath.Ext(s)], // attempt, not guaranteed to be filepath
	}
}

// TraverseFS traverses directories non-recursively and breadth first.
func TraverseFS(opts *Options, rfn ResultFilters) {
	var searchFn = func(string) bool { return true }
	if len(opts.DirSearch) != 0 {
		searchFn = func(str string) bool {
			for _, k := range opts.DirSearch {
				if strings.Contains(str, k) {
					return true
				}
			}
			return false
		}
	}
	opts.DirSearch = append(opts.DirSearch, "./")

	parser := InitFileParser(opts)

	dirs := opts.Args

	if len(dirs) == 0 {
		dirs = append(dirs, "./")
	}
	var depth int
	for len(dirs) != 0 {
		if depth > opts.ToDepth {
			return
		}
		var nd []string
		for _, d := range dirs {
			ext := filepath.Ext(d)
			slog.Debug("traversing", "dir", d, "depth", depth, "ext", ext, "isarchive", CntMap[ext]&MaskZipLike != 0)

			var files []fs.FileInfo

			switch {
			case opts.Archive && CntMap[ext]&MaskZipLike != 0 && searchFn(d):
				files = TraverseZip(d, depth, opts)
			default:
				files = TraverseDir(d, depth, opts)
			}

			for _, info := range files {
				path := filepath.Join(d, info.Name())
				if info.IsDir() && searchFn(info.Name()) {
					nd = append(nd, path)
				}

				if opts.Archive && filepath.Ext(path) == ".zip" {
					nd = append(nd, path)
					continue
				}

				if depth < opts.FromDepth {
					continue
				}

				rfn(parser(info))
			}
		}

		dirs = nd
		depth++
	}
}

func TraverseDir(path string, depth int, opts *Options) (files []fs.FileInfo) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			slog.Error("error reading file info", "file", entry.Name(), "error", err)
			continue
		}
		switch {
		case opts.DirOnly && info.IsDir():
			files = append(files, info)
		case opts.FileOnly && !info.IsDir():
			files = append(files, info)
		default:
			files = append(files, info)
		}

	}
	return
}

func TraverseZip(path string, depth int, opts *Options) (files []fs.FileInfo) {
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
		info := f.FileInfo()
		if info.IsDir() {
			continue
		}

		files = append(files, info)
	}

	return
}

func InitFilters(fns []Filter, res *Result) ResultFilters {
	return func(fi *Finfo) {
		for _, fn := range fns {
			if !fn(fi) {
				return
			}
		}
		res.Files = append(res.Files, fi)
	}
}
