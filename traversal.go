package main

import (
	"archive/zip"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Traverse traverses directories non-recursively and breadth first.
func Traverse(wfn fs.WalkDirFunc) {
	dirs := args
	var depth int
	for len(dirs) != 0 {
		if depth > Opts.ToDepth {
			return
		}
		var nd []string
		for _, d := range dirs {

			if filepath.Ext(d) == ".zip" {
				traverseZip(d, depth, wfn)
				continue
			}

			entries, err := os.ReadDir(d)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				path := filepath.Join(d, entry.Name())
				if entry.IsDir() {
					nd = append(nd, path)
				}

				if Opts.Archive && filepath.Ext(path) == ".zip" {
					nd = append(nd, path)
					continue
				}

				if depth < Opts.FromDepth {
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

func traverseZip(path string, depth int, wfn fs.WalkDirFunc) {
	r, err := zip.OpenReader(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Close()

	for _, f := range r.File {
		fn := filepath.ToSlash(f.Name)

		fdepth := depth + strings.Count(fn, "/")
		if fdepth < Opts.FromDepth || fdepth > Opts.ToDepth {
			continue
		}
		fpath := filepath.Join(path, fn)

		entry := ZipEntry{f}

		if depth <= Opts.FromDepth {
			continue
		}

		err := wfn(fpath, entry, nil)
		if err != nil {
			continue
		}
	}
}

func buildWalkDirFn(fns []filter, res *result) func(string, fs.DirEntry, error) error {
	return func(path string, d fs.DirEntry, err error) error {
		if d == nil || err != nil {
			return nil
		}
		fi := &finfo{name: d.Name(), path: path}
		for _, fn := range fns {
			res := fn(fi, d)
			if !res {
				return nil
			}
		}
		res.files = append(res.files, fi)
		return nil
	}
}
