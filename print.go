package main

import (
	"bufio"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
)

// the interval of flushing the buffer
const bufLength = 500

func printWithBuf(files []*finfo) {
	if len(files) == 0 {
		return
	}
	if Opts.Quiet {
		slog.Debug("quiet flag is set, returning from print function")
		return
	}

	if Opts.Tree {
		ftree := AddFilesToTree(files)
		ftree.PrintTree("")
		return
	}

	// I am unsure of how large this buffer should be. Testing or profiling might be necessary to
	// find what is reasonable. The default buffer size was flushing automatically before being told to.
	// This might be okay in itself, and we might not need to manually set a buffer ta all (or flush).
	var str []byte

	w := bufio.NewWriterSize(os.Stdout, 4096*bufLength)

	for i, file := range files {
		fp := filepath.ToSlash(file.path)
		if Opts.Absolute {
			fp, _ = filepath.Abs(file.path)
			fp = filepath.ToSlash(fp)
		}
		res := fp + "\n"

		if Opts.Clipboard {
			str = append(str, res...)
		}

		w.WriteString(res)
		if i%bufLength == 0 {
			w.Flush()
		}
	}

	w.Flush()

	if Opts.Clipboard {
		if str[len(str)-1] == '\n' {
			str = str[:len(str)-1]
		}
		clipboard.WriteAll(string(str))
	}
}
