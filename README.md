# list: Filesystem listing simplified.
## Overview 
`list` is not a replacement for the standard `ls` command. It is designed as a simplified CLI tool for specific operations concerning file and directory listing, with an added layer of pattern matching functionality. Files are ordered by name in descending order by default, with options for ordering by modification time and ascending order.

*Especially for PowerShell, where the `ls` output may not be pipeable to certain programs by default.*
## Features
- List files in a directory, showing their relative locations.
- Recursively list files in subdirectories.
- Include and exclude files based on pre-defined extension patterns.
- Sort files by either name or modified time, in ascending or descending order.
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

5. (Optional) Rename binary `main` or `main.exe` into `list` or `list.exe` and include the binary in your `path`.

## Command-line Flags

|Flag|Short|Long|Description|
|---|---|---|---|
|Absolute|`-A`|`--absolute`|Format paths to be absolute. Relative by default.|
|Recurse|`-r`|`--recurse`|Recursively list files in subdirectories|
|Exclusive recurse|`-x`|`--xrecurse`|Exclusively list files in subdirectories|
|Ascending|`-a`|`--ascending`|Order results in ascending order. Descending by default.|
|Date|`-d`|`--date`|Order results by modified time. Ordered by filename by default.|
|Query|`-q`|`--query`|Return the most similar file to the query.|
|Include|`-i`|`--include`|Include only items fitting a given extension pattern. Use ',' for multiple patterns.|
|Exclude|`-e`|`--exclude`|Exclude items fitting a given extension pattern. Use ',' for multiple patterns.|


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