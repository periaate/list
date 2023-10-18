# ls-go: Filesystem listing simplified.
## Overview 
`ls-go` is not a replacement for the standard `ls` command. It is designed as a simplified CLI tool for specific operations concerning file and directory listing, with an added layer of pattern matching functionality. Files are ordered by name in descending order by default, with options for ordering by modification time and ascending order.

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
git clone https://github.com/periaate/ls-go
```

3. `cd` into the project folder:
```bash
cd ls-go
```

4. Build the project:
```bash
go build main.go
```

5. (Optional) Rename binary `main` or `main.exe` into `ls-go` or `ls-go.exe` and include the binary in your `path`.

## Command-line Flags

|Flag|Short|Long|Description|
|---|---|---|---|
|Recurse|`-r`|`--recurse`|Recursively list files in subdirectories|
|Exclusive recurse|`-x`|`--xrecurse`|Exclusively list files in subdirectories|
|Ascending|`-A`|`--ascending`|Results will be ordered in ascending order. Descending by default.|
|Date|`-d`|`--date`|Results will be ordered by modified time. Ordered by filename by default.|
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
ls-go
```

2. Recursively list all files.
```bash
ls-go -r
```

3. List files in ascending order.
```bash
ls-go -A
```

4. List files sorted by modified date.
```bash
ls-go -d
```

5. Include only image and document files.
```bash
ls-go -i image,document
```

6. Exclude video files.
```bash
ls-go -e video
```
