package list

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/periaate/common"
)

type Filter func(*Element) bool

const (
	Other    = "other"
	Image    = "image"
	Video    = "video"
	Audio    = "audio"
	Media    = "media"
	Archive  = "archive"
	ZipLike  = "zip"
	Code     = "code"
	Conf     = "conf"
	Docs     = "docs"
	OtherDev = "odev"

	_ uint32 = 1 << iota
	MaskImage
	MaskVideo
	MaskAudio
	MaskArchive
	MaskZipLike = 1<<iota + MaskArchive
	MaskCode    = 1 << iota
	MaskConf
	MaskDocs
	MaskOtherDev
)

// Hide contains commonly unwanted files and directories. Any beginning with a dot hidden by default.
var Hide = map[string]bool{
	"Thumbs.db":                 true,
	"desktop.ini":               true,
	"Icon\r":                    true,
	"System Volume Information": true,
	"$RECYCLE.BIN":              true,
	"lost+found":                true,
	"node_modules":              true,
}

var CntMasks = map[uint32][]string{
	MaskImage:    {".jpg", ".jpeg", ".png", ".apng", ".gif", ".bmp", ".webp", ".avif", ".jxl", ".tiff"},
	MaskVideo:    {".mp4", ".m4v", ".webm", ".mkv", ".avi", ".mov", ".mpg", ".mpeg"},
	MaskAudio:    {".m4a", ".opus", ".ogg", ".mp3", ".flac", ".wav", ".aac"},
	MaskArchive:  {".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".lz4", ".zst", ".lzma", ".lzip", ".lz", ".cbz"},
	MaskZipLike:  {".zip", ".cbz", ".cbr"},
	MaskCode:     {".go", ".c", ".h", ".cpp", ".hpp", ".rs", ".py", ".js", ".ts", ".html", ".css", ".scss", ".java", ".php"},
	MaskConf:     {".json", ".toml", ".yaml", ".yml", ".xml", ".ini", ".cfg", ".conf", ".properties", ".env"},
	MaskDocs:     {".pdf", ".epub", ".mobi", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".odt", ".ods", ".odp", ".txt", ".rtf", ".csv", ".tsv", ".md"},
	MaskOtherDev: {".sql", ".sh", ".bat", ".cmd", ".ps1", ".psm1", ".psd1", ".ps1xml", ".pssc", ".psc1", ".pssc", ".psh"},
}
var CntMap = map[string]uint32{}

func AsMask(sar []string) uint32 {
	var mask uint32
	for _, v := range sar {
		mask |= StrToMask(v)
	}
	return mask
}

func StrToMask(str string) uint32 {
	switch str {
	case Image, "img", "i":
		return MaskImage
	case Video, "vid", "v":
		return MaskVideo
	case Audio, "a":
		return MaskAudio
	case Media, "m":
		return MaskImage | MaskVideo | MaskAudio
	case Archive:
		return MaskArchive
	case ZipLike:
		return MaskZipLike
	case Code:
		return MaskCode
	case Conf:
		return MaskConf
	case Docs:
		return MaskDocs
	case OtherDev:
		return MaskOtherDev
	default:
		return 0
	}
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

func CollectFilters(opts *Options) Filter {
	switch {
	case opts.DirOnly:
		opts.Filters = append(opts.Filters, func(fi *Element) bool {
			return fi.IsDir
		})
	case opts.FileOnly:
		opts.Filters = append(opts.Filters, func(fi *Element) bool {
			return !fi.IsDir
		})
	}

	if (len(opts.Search)) > 0 {
		opts.Filters = append(opts.Filters, ParseSearch(opts.Search))
	}

	return common.All(true, opts.Filters...)
}

func ParseSearch(args []string) Filter {
	slog.Debug("search args", "args", args)
	filters := []func(*Element) bool{}

	for _, arg := range args {
		q := Query{Include: true}
		switch {
		case len(arg) < 2:
			continue
		case arg[:2] == "-=":
			arg = arg[1:]
			q.Include = false
			fallthrough
		case arg[0] == '=':
			q.Value = arg
			q.Kind = Exact
		case arg[0] == '-':
			arg = arg[1:]
			q.Include = false
			fallthrough
		default:
			q.Kind = Substring
			q.Value = arg
		}
		filters = append(filters, q.GetFilter())
	}

	return common.All(true, filters...)
}

func ParseKind(args []string, inc bool) Filter {
	q := Query{
		Kind:    MaskK,
		Include: inc,
	}
	for _, arg := range args {
		q.Mask |= CntMap[filepath.Ext(arg)]
	}
	return q.GetFilter()
}

type QueryKind [2]bool

var (
	Substring = QueryKind{false, false}
	Fuzzy     = QueryKind{true, false}
	Exact     = QueryKind{false, true}
	MaskK     = QueryKind{true, true}
)

type Query struct {
	Value   string
	Include bool
	Mask    uint32
	Kind    QueryKind
}

func (q Query) GetFilter() (f Filter) {
	switch q.Kind {
	case Fuzzy:
		f = QueryAsFilter(q.Value)
	case Exact:
		f = ExactFilter(q.Value)
	case Substring:
		fallthrough
	default:
		f = SubstringFilter(q.Value)
	}

	if !q.Include {
		f = common.Negate(f)
	}
	return f
}

func ExactFilter(search string) Filter {
	return func(e *Element) bool { return search == e.Name }
}
func SubstringFilter(search string) Filter {
	return func(e *Element) bool {
		r := strings.Contains(e.Name, search)
		slog.Debug("substring filter", "name", e.Name, "search", search, "result", r)
		return r
	}
}
