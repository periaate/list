# list
## Overview 
`list` is a directory traversal program with built in searching, filtering, and sorting functionalities. 

## Usage:
`list [OPTIONS]`

### Options
Traversal and filtering options are called at the same time. Then Processing options, then printing options. The order in which options are listed in this document mirrors the order in which they are evaluated or utilized.

#### Traversal options
Determines how the traversal is done.:\
  `-r`, `--recurse`     Recursively list files in subdirectories. Directory traversal is done iteratively and breadth first.\
  `-z`                  Treat zip archives as directories.\
  `-T`, `--todepth=`    List files to a certain depth. (default: 0)\
  `-F`, `--fromdepth=`  List files from a certain depth. (default: -1)

#### Filtering options
Applied while traversing, called on every entry found.:\
  `-s`, `--search=`     Only include items which have search terms as substrings. Can be used multiple times. Multiple values are inclusive by default. (OR)\
        `--AND`         Including this flag makes search with multiple values conjuctive, i.e., all search terms must be matched. (AND)\
  `-i`, `--include=`    File type inclusion. Can be used multiple times.\
  `-e`, `--exclude=`    File type exclusion. Can be used multiple times.\
  `[image|video|audio|archive|ziplike]`\
  `-I`, `--ignore=`     Ignores all paths which include any given strings.\
        `--dirs`        Only include directories in the result.\
        `--files`       Only include files in the result.

#### Processing options
Applied after traversal, called on the final list of files.:\
  `-q`, `--query=`      Fuzzy search query. Results will be ordered by their score.\
  `-a`, `--ascending`   Results will be ordered in ascending order.\
  `-S`, `--sort=`       Sort the result by word. (default: none)\
      `[none|name|n|mod|time|t|size|s|creation|c]`\
      `--select=`     Select a single element or a range of elements. Usage: `[{index}]` `[{from}:{to}]` `[{from}:{to}={page}]` Supports negative indexing, and relative `+` indexing. Can be used without a flag as the last argument.\
      `--shuffle`     Randomly shuffle the result.\
      `--seed=`       Seed for the random shuffle. (default: -1)

#### Printing options
Determines how the results are printed.:\
  `-A`, `--absolute`    Format paths to be absolute. Relative by default.\
  `-D`, `--debug`       Debug flag enables debug logging.\
  `-Q`, `--quiet`       Quiet flag disables printing results.\
  `-c`, `--clipboard`   Copy the result to the clipboard.\
        `--tree`        Prints as tree.

### Examples
All of the examples implicitly traverse from the current working directory `./`.
Traverse recursively and list last 10 results:\
`list -r [-10:]`

Traverse recursively, listing only the files in the immediate subdirectories, but not their subdirectories:\
`list -T 2 -F 1`

Traverse recursively, ignoring all directories with the substrings `".git"` and `".thumbs"`, only including image files, only including files with the substring `"_p01"`, traversing archives as directories, fuzzy searching with queries `"picasso"` and `"museum"`  and sorting by score, and printing only the top 100 files with absolute paths:\
`list -rAz -I ".git" -I ".thumbs" -s "_p01" -q "picasso" -q "museum" -i image [:100]`
