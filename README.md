# list: Filesystem listing program
## Overview 
`list` is a command-line utility that offers extended functionality for file and directory listing, allowing for some complex behavior while remaining simple to use.

## Help
```
Usage:
  list [OPTIONS] [SELECTION]

Traversal options - Determines how the traversal is done.:
  -r, --recurse     Recursively list files in subdirectories
  -z                Treat zip archives as directories.
  -T, --todepth:    List files to a certain depth. (default: 0)
  -F, --fromdepth:  List files from a certain depth. (default: -1)

Filtering options - Applied while traversing, called on every entry found.:
  -i, --include:    File type inclusion: image, video, audio
  -e, --exclude:    File type exclusion: image, video, audio.
  -I, --ignore:     Ignores all paths which include any given strings.
  -s, --search:     Only include paths which include any given strings.
      --dirs        Only include directories in the result.
      --files       Only include files in the result.

Processing options - Applied after traversal, called on the final list of files.:
  -q, --query:      Fuzzy search query. Results will be ordered by their score.
  -a, --ascending   Results will be ordered in ascending order. Files are ordered into descending order by default.
  -d, --date        Results will be ordered by their modified time. Files are ordered by filename by default
  -S, --select:     Select a single element or a range of elements. Usage: [{index}] [{from}:{to}] Supports negative
                    indexing. Can be used without a flag as the last argument.
  -n, --sort        Sort the result. Files are ordered by filename by default.

Printing options - Determines how the results are printed.:
  -A, --absolute    Format paths to be absolute. Relative by default.
  -D, --debug       Debug flag enables debug logging.
  -Q, --quiet       Quiet flag disables printing results.

Help Options:
  -?                Show this help message
  -h, --help        Show this help message
```
