package main

import (
	"bufio"
	"log/slog"
	"os"
	"path/filepath"
)

// the interval of flushing the buffer
const bufLength = 500

func printWithBuf(files []*finfo) {
	if Opts.Quiet {
		slog.Debug("quiet flag is set, returning from print function")
		return
	}
	// I am unsure of how large this buffer should be. Testing or profiling might be necessary to
	// find what is reasonable. The default buffer size was flushing automatically before being told to.
	// This might be okay in itself, and we might not need to manually set a buffer ta all (or flush).
	buf := bufio.NewWriterSize(os.Stdout, 4096*bufLength)

	for i, file := range files {
		fp := filepath.ToSlash(file.path)
		if Opts.Absolute {
			fp, _ = filepath.Abs(file.path)
			fp = filepath.ToSlash(fp)
		}

		buf.WriteString(fp + "\n")
		if i%bufLength == 0 {
			buf.Flush()
		}
	}
	buf.Flush()
}
