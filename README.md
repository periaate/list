# list: Filesystem listing program
## Overview 
`list` is a command-line utility that offers extended functionality for file and directory listing, allowing for some complex behavior while remaining simple to use.

## Options
```
  -A, --absolute    Format paths to be absolute. Relative by default.
  -r, --recurse     Recursively list files in subdirectories
  -T, --todepth:    List files to a certain depth.
  -F, --fromdepth:  List files from a certain depth.
  -i, --include:    File type inclusion: image, video, audio
  -e, --exclude:    File type exclusion: image, video, audio.
  -I, --ignore:     Ignores all paths which include any given strings.
  -s, --search:     Only include paths which include any given strings.
      --dirs        Only include directories in the result.
      --files       Only include files in the result.
  -q, --query:      Fuzzy search query. Results will be ordered by their score.
  -a, --ascending   Results will be ordered in ascending order. Files are ordered into descending order by default.
  -d, --date        Results will be ordered by their modified time. Files are ordered by filename by default
  -S, --slice:      Slice [{from}:{to}]. Supports negative indexing.
```

## Note
This project was recently refactored, out of date documentation removed and new documentation will be added later.