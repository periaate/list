package list

import (
	"log/slog"
	"os"
	"os/exec"
)

func Exec(res *Result, opts *Options) {
	opts.ExecArgs = append(opts.ExecArgs, res.Sar()...)

	cmd := exec.Command(opts.ExecArgs[0], opts.ExecArgs[1:]...)
	cmd.Stdout = os.Stdout
	if opts.Debug {
		cmd.Stderr = os.Stderr
	}
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		slog.Debug("error running command", "err", err)
	}
}
