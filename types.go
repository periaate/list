package main

import "path/filepath"

type contentType uint8

func getContentType(filename string) contentType {
	ext := filepath.Ext(filename)
	if t, ok := contentTypes[ext]; ok {
		return t
	}
	return other
}

func stringToContentType(s string) contentType {
	switch s {
	case "image":
		return image
	case "video":
		return video
	case "audio":
		return audio
	default:
		return other
	}
}

const (
	other contentType = (1 << iota) - 1
	image
	video
	audio
)

var contentTypes = map[string]contentType{
	// image
	".jpg":  image,
	".jpeg": image,
	".png":  image,
	".apng": image,
	".gif":  image,
	".bmp":  image,
	".webp": image,
	".avif": image,
	".jxl":  image,
	".tiff": image,

	// video
	".mp4":  video,
	".m4v":  video,
	".webm": video,
	".mkv":  video,
	".avi":  video,
	".mov":  video,
	".mpg":  video,
	".mpeg": video,

	// audio
	".m4a":  audio,
	".opus": audio,
	".ogg":  audio,
	".mp3":  audio,
	".flac": audio,
	".wav":  audio,
	".aac":  audio,
}
