package list

import (
	"bufio"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
)

// the interval of flushing the buffer
const bufLength = 500

func PrintWithBuf(files []*Finfo, opts *Options) {
	if len(files) == 0 {
		return
	}
	if opts.Quiet {
		slog.Debug("quiet flag is set, returning from print function")
		return
	}

	if opts.Tree {
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
		fp := filepath.ToSlash(file.Path)
		if opts.Absolute {
			fp, _ = filepath.Abs(file.Path)
			fp = filepath.ToSlash(fp)
		}
		res := fp + "\n"

		if opts.Clipboard {
			str = append(str, res...)
		}

		w.WriteString(res)
		if i%bufLength == 0 {
			w.Flush()
		}
	}

	w.Flush()

	if opts.Clipboard {
		if str[len(str)-1] == '\n' {
			str = str[:len(str)-1]
		}
		clipboard.WriteAll(string(str))
	}
}
