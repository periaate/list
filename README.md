# list: Filesystem listing program
## Overview 
`list` is a command-line utility that offers extended functionality for file and directory listing, allowing for some complex behavior while remaining simple to use. It primarily enhances file system navigation with features file filtering based on extension patterns or built-in fuzzy string matching and sorting.

## Features
- **File Listing**: Displays files in directories, with options for absolute or relative path formats.
- **Recursive Directory Listing**: Lists files in subdirectories, with an option for exclusive recursion into subdirectories.
- **Sorting**: Sort files by name or modified time, in either ascending or descending order.
- **Pattern Matching**: Includes or excludes files based on specified extension patterns.
- **Fuzzy String Matching**: Offers fuzzy string matching, either listing items sorted by how closely they match the query, or selecting the item which matched the closest to the query.

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

|Flag|Short|Long|Description|
|---|---|---|---|
|Absolute|`-A`|`-absolute`|Format paths to be absolute. Relative by default.|
|Recurse|`-r`|`-recurse`|Recursively list files in subdirectories|
|Exclusive Recurse|`-x`|`-xrecurse`|Exclusively list files in subdirectories|
|Ascending|`-a`|`-ascending`|Results will be ordered in ascending order. Files are ordered into descending order by default.|
|Date|`-d`|`-date`|Results will be ordered by their modified time. Files are ordered by filename by default|
|Select|`-S`|`-select`|Selects the item which matches the query the best.|
|Query|`-Q`|`-query`|Returns items ordered by their similarity to the query.|
|Include|`-i`|`-include`|Given an existing extension pattern configuration target, will include only items fitting the pattern. Use ',' to define multiple patterns.|
|Exclude|`-e`|`-exclude`|Given an existing extension pattern configuration target, will exclude items fitting the pattern. Use ',' to define multiple patterns.|


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

## Usage Examples

1. List all files in the current directory.
```bash
list
```

2. Recursively list all files.
```bash
list -r
```

3. List files in ascending order.
```bash
list -a
```

4. List files sorted by modified date.
```bash
list -d
```

5. Include only image and document files.
```bash
list -i image,document
```

6. Exclude video files.
```bash
list -e video
```

7. List with absolute paths.
```bash
list -A
```

8. Only list files in subdirectories.
```bash
list -x
```