# list: Filesystem listing program
## Overview 
`list` is a command-line utility that offers extended functionality for file and directory listing, allowing for some complex behavior while remaining simple to use.

## Features
- **File Listing**: Displays files in directories, with options for absolute or relative path formats.
- **Recursive Directory Listing**: Lists files in subdirectories, with an option for exclusive recursion into subdirectories.
- **Sorting**: Sort files by name or modified time, in either ascending or descending order.
- **Pattern Matching**: Includes or excludes files based on specified extension patterns.
- **Fuzzy String Matching**: Performant and configurable string matching using n-grams.

## Building from source
1. Traverse to `src` under the `$GOPATH`.
For Linux:
```bash
cd ~/go/src/
```
For PowerShell:
```PowerShell
cd $HOME\go\src\
```

2. Clone the repository there:
```bash
git clone https://github.com/periaate/list
```

3. `cd` into the project folder:
```bash
cd list
```

4. Build the project:
```bash
go build main.go
```

5. (Optional) Rename binary `main` or `main.exe` into `list` or `list.exe` (or any other name) and include the binary in your `path`.

## Command-line Flags

| Flag               | Short | Long        | Description |
|--------------------|-------|-------------|-------------|
| Absolute           | `-A`  | `--absolute` | Format paths to be absolute. Relative by default. |
| Recurse            | `-r`  | `--recurse`  | Recursively list files in subdirectories |
| Exclusive Recurse  | `-x`  | `--xrecurse` | Exclusively list files in subdirectories |
| Invert             | `-a`  | `--invert`   | Results will be ordered in inverted order. Files are ordered into descending order by default. |
| Date               | `-d`  | `--date`     | Results will be ordered by their modified time. Files are ordered by filename by default. |
| Combine            | `-c`  | `--combine`  | If given multiple paths, will combine the results into one list. |
| Query              | `-q`  | `--query`    | Returns items ordered by their similarity to the query. |
| Ngram              | `-n`  | `--ngram`    | Defines the n-gram size for the query algorithm. Default is 3. |
| Top                | `-t`  | `--top`      | Returns first n items. |
| Prune              | `-p`  | `--prune`    | Prunes items with a score lower than the given value. -1 to prune all items without a score. 0.0 is the default and will not prune any items. |
| Score              | `-s`  | `--score`    | Returns items with their score. Intended for debugging purposes. |
| Include            | `-i`  | `--include`  | Given an existing extension pattern configuration target, will include only items fitting the pattern. Use ',' to define multiple patterns. |
| Exclude            | `-e`  | `--exclude`  | Given an existing extension pattern configuration target, will exclude items fitting the pattern. Use ',' to define multiple patterns. |



## Pattern File (`patterns.yml`)

The `patterns.yml` file contains predefined file extension patterns, categorized under various types like `image`, `video`, `audio`, `data`, and `document`.

Here is an excerpt from the default `patterns.yml`.
```yml
extensions:
    image:
        - .jpg
        - .png
        - # ... other image extensions
    # ... other categories
```

You may modify or add patterns into it which will be included on recompilation. There is currently no way to have dynamically configured patterns, although this is a planned feature.
