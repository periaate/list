# list: Filesystem listing program
## Overview 
`list` is a directory traversal program with built in searching, filtering, and sorting functionalities. 

## Usage:
`list [OPTIONS] [SELECTION]`

### Options
Traversal and filtering options are called at the same time. Then Processing options, then printing options. The order in which options are listed in this document mirrors the order in which they are evaluated or utilized.

#### Traversal options
Determines how the traversal is done.:\
  `-r`, `--recurse`     Recursively list files in subdirectories. Directory traversal is done iteratively and breadth first.\
  `-z`                  Treat zip archives as directories.\
  `-T`, `--todepth`:    List files to a certain depth.\
  `-F`, `--fromdepth`:  List files from a certain depth.

#### Filtering options
Applied while traversing, called on every entry found.:\
  `-s`, `--search`:     Only include items which have search terms as substrings. Can be used multiple times. Multiple values are inclusive by default. (OR)\
  `-S`, `--xsearch`:    Including this flag makes search with multiple values conjuctive. (AND)
  `-i`, `--include`:    File type inclusion: image, video, audio. Can be used multiple times.\
  `-e`, `--exclude`:    File type exclusion: image, video, audio. Can be used multiple times.\
  `-I`, `--ignore`:     Ignores all paths which include any given strings. Can be used multiple times.\
      `--dirs`        Only include directories in the result.\
      `--files`       Only include files in the result.

#### Processing options
Applied after traversal, called on the final list of files.:\
  `-q`, `--query`:      Fuzzy search query. Results will be ordered by their score. Can be used multiple times.\
  `-a`, `--ascending`   Results will be ordered in ascending order. Files are ordered into descending order by default.\
  `-d`, `--date`        Results will be ordered by their modified time. Files are ordered by filename by default.\
  `-n`, `--sort`        Sort the result. Files are ordered by filename by default.\
  `-S`, `--select`:     Select a single element or a range of elements. Usage: [{index}] or [{from}:{to}] Supports negative indexing. Can be used without a flag as the last argument.

#### Printing options
Determines how the results are printed.:\
  `-A`, `--absolute`    Format paths to be absolute. Relative by default.\
  `-D`, `--debug`       Debug flag enables debug logging.\
  `-Q`, `--quiet`       Quiet flag disables printing results.

### Examples
All of the examples implicitly traverse from the current working directory `./`.
Traverse recursively and list last 10 results:\
`list -r [-10:]`

Traverse recursively, listing only the files in the immediate subdirectories, but not their subdirectories:\
`list -T 2 -F 1`

Traverse recursively, ignoring all directories with the substrings `".git"` and `".thumbs"`, only including image files, only including files with the substring `"_p01"`, traversing archives as directories, fuzzy searching with queries `"picasso"` and `"museum"`  and sorting by score, and printing only the top 100 files with absolute paths:\
`list -rAz -I ".git" -I ".thumbs" -s "_p01" -q "picasso" -q "museum" -i image [:100]`
