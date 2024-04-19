package ls

import (
	"math"

	"github.com/periaate/list"
	"github.com/periaate/slice"
)

const (
	Images   = list.Image
	Videos   = list.Video
	Audios   = list.Audio
	Medias   = list.Media
	Docs     = list.Docs
	Conf     = list.Conf
	Code     = list.Code
	Archives = list.Archive
	ZipLikes = list.ZipLike
)

type OptFn func(*list.Options)

func Paths(paths ...string) OptFn {
	return func(opts *list.Options) { opts.Args = append(opts.Args, paths...) }
}

func FromDepth(depth int) OptFn {
	return func(opts *list.Options) { opts.FromDepth = depth }
}

func ToDepth(depth int) OptFn {
	return func(opts *list.Options) { opts.ToDepth = depth }
}

func DepthPattern(pat string) OptFn {
	return func(opts *list.Options) {
		from, to, err := slice.ParsePattern(pat, math.MaxInt)
		if err != nil {
			return
		}
		opts.FromDepth = from
		opts.ToDepth = to
	}
}

const (
	Files = false
	Dirs  = true
)

func Only(b bool) OptFn {
	return func(opts *list.Options) {
		switch b {
		case Files:
			opts.OnlyFiles = true
		case Dirs:
			opts.DirOnly = true
		}
	}
}

func Search(queries ...string) OptFn {
	return func(opts *list.Options) {
		f := list.ParseSearch(queries)
		if f != nil {
			opts.Filters = append(opts.Filters, f)
		}
	}
}

const (
	Include = true
	Exclude = false
)

func Kind(inc bool, kinds ...string) OptFn {
	return func(opts *list.Options) {
		f := list.ParseKind(kinds, inc)
		if f != nil {
			opts.Filters = append(opts.Filters, f)
		}
	}
}

func Recurse(opts *list.Options) {
	opts.ToDepth = math.MaxInt64
}

func NoHide(opt *list.Options) {
	opt.NoHide = true
}

func Do(args ...OptFn) []*list.Element {
	opts := &list.Options{}
	for _, fn := range args {
		fn(opts)
	}

	if len(opts.Args) == 0 {
		opts.Args = []string{"./"}
	}

	return list.Run(opts)
}
