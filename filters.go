package list

import (
	"strings"
)

type Filter func(*Finfo) bool

const (
	Other   = "other"
	Image   = "image"
	Video   = "video"
	Audio   = "audio"
	Archive = "archive"
	ZipLike = "zip"
	Media   = "media"

	_ uint32 = 1 << iota
	MaskImage
	MaskVideo
	MaskAudio
	MaskArchive
	MaskZipLike = 1<<iota + MaskArchive
	// MaskMedia   = MaskImage | MaskVideo | MaskAudio
)

var CntMasks = map[uint32][]string{
	MaskImage:   {".jpg", ".jpeg", ".png", ".apng", ".gif", ".bmp", ".webp", ".avif", ".jxl", ".tiff"},
	MaskVideo:   {".mp4", ".m4v", ".webm", ".mkv", ".avi", ".mov", ".mpg", ".mpeg"},
	MaskAudio:   {".m4a", ".opus", ".ogg", ".mp3", ".flac", ".wav", ".aac"},
	MaskArchive: {".zip", ".rar", ".7z", ".tar", ".gz", ".bz2", ".xz", ".lz4", ".zst", ".lzma", ".lzip", ".lz", ".cbz"},
	MaskZipLike: {".zip", ".cbz", ".cbr"},
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
	case Image:
		return MaskImage
	case Video:
		return MaskVideo
	case Audio:
		return MaskAudio
	case Archive:
		return MaskArchive
	case ZipLike:
		return MaskZipLike
	// case Media:
	// 	return MaskMedia
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

	// var med []string
	// med = append(med, CntMasks[MaskImage]...)
	// med = append(med, CntMasks[MaskVideo]...)
	// med = append(med, CntMasks[MaskAudio]...)
	// RegisterMasks(MaskMedia, med...)
}

func CollectFilters(opts *Options) []Filter {
	var fns []Filter
	switch {
	case opts.DirOnly:
		fns = append(fns, func(fi *Finfo) bool {
			return fi.IsDir
		})
	case opts.FileOnly:
		fns = append(fns, func(fi *Finfo) bool {
			return !fi.IsDir
		})
	}

	if (len(opts.Search) + len(opts.Include) + len(opts.Exclude) + len(opts.Ignore)) > 0 {
		fns = append(fns, FilterList(opts))
	}
	return fns
}

func FilterList(opts *Options) Filter {
	incMask := AsMask(opts.Include)
	excMask := AsMask(opts.Exclude)

	var searchFn func(string) bool
	if opts.SearchAnd {
		searchFn = func(str string) bool {
			for _, sub := range opts.Search {
				if !strings.Contains(str, sub) {
					return false
				}
			}
			return true
		}
	} else {
		searchFn = func(str string) bool {
			for _, sub := range opts.Search {
				if strings.Contains(str, sub) {
					return true
				}
			}
			return false
		}
	}

	return func(fi *Finfo) bool {
		any := searchFn(fi.Name)
		if len(opts.Search) > 0 && !any {
			return false
		}

		if incMask > 0 && (incMask&fi.Mask) == 0 {
			return false
		}

		for _, ign := range opts.Ignore {
			if strings.Contains(fi.Path, ign) {
				return false
			}
		}

		if excMask > 0 && (excMask&fi.Mask) != 0 {
			return false
		}

		return true
	}
}
