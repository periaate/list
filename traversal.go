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

		// fmt.Println(
		// 	fi.Mask&MaskImage != 0,
		// 	fi.Mask&MaskVideo != 0,
		// 	fi.Mask&MaskAudio != 0,
		// 	fi.Mask&MaskArchive != 0,
		// 	fi.Mask&MaskZipLike != 0,
		// )

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
