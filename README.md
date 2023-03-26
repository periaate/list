# ls-go
A 🚀 Blazing Fast 🚀 Directory Listing Utility for the Command Line

## Overview
Introducing the Blazing Fast Directory Lister, a high-performance, user-friendly, and cross-platform command-line utility for effortlessly listing files and directories on your computer! Built with cutting-edge technology, this utility is designed to make your life easier and more productive! 💪

## Installation
The Blazing Fast Directory Lister is available for all major platforms, including Windows, macOS, and Linux. Follow the step-by-step installation examples below for your platform:

### Windows
Install the Go compiler: Download and install the latest version of Go from the official website here.

Clone the Blazing Fast Directory Lister repository or download the source code as a ZIP file and extract it to your preferred location.

Open a command prompt (cmd) or PowerShell and navigate to the directory containing the source code.

Compile the program by running the following command:

`go build -o ls-go.exe`

Add the compiled executable to your system PATH:

**For cmd:**
a. Open the "Environment Variables" settings by searching for "Edit the system environment variables" in the Start menu.

b. In the "System Variables" section, find the variable named "Path", select it, and click the "Edit..." button.

c. Click the "New" button and add the path to the directory containing the `ls-go.exe` executable.

**For PowerShell:**
a. Open PowerShell and run the following command (replacing "path\to\ls-go" with the actual path to the directory containing the `ls-go.exe` executable):

`[Environment]::SetEnvironmentVariable("Path", $env:Path + ";path\to\ls-go", "User")`
Close and reopen your command prompt or PowerShell to apply the changes.

You can now use the Blazing Fast Directory lister by typing `ls-go` in your command prompt or PowerShell!

### macOS and Linux
Install the Go compiler: Download and install the latest version of Go from the official website [here](https://go.dev/doc/install).

Clone the Blazing Fast Directory Lister repository or download the source code as a ZIP file and extract it to your preferred location.

Open a terminal and navigate to the directory containing the source code.

Compile the program by running the following command:

`go build -o ls-go`
Add the compiled executable to your system PATH:

**a.** Open your shell configuration file (e.g., `~/.bashrc` for bash or `~/.zshrc` for zsh) in a text editor.

**b.** Add the following line at the end of the file (replacing "path/to/ls-go" with the actual path to the directory containing the ls-go executable):

`export PATH=$PATH:path/to/ls-go`
**c.** Save the file and close the text editor.

Restart your terminal or run `source ~/.bashrc` (for bash) or `source ~/.zshrc` (for zsh) to apply the changes.

You can now use the Blazing Fast Directory Lister by typing `ls-go` in your terminal!

## Usage
The Blazing Fast Directory Lister is incredibly simple to use! Just follow the examples below to get started:

To list the files in the current directory:

`ls-go`

To list the files in a specific directory:

`ls-go /path/to/directory`


- To list the files in the current directory and its subdirectories (recursively):

`ls-go -r`


- To list the files in a specific directory and its subdirectories (recursively):

`ls-go /path/to/directory -r`


- To list the files in a specific directory and its subdirectories (recursively) using alternative flag options:

`ls-go /path/to/directory --recurse`

or

`ls-go /path/to/directory -R`


Feel free to mix and match the directory paths and flags as needed. The Blazing Fast Directory Lister will handle it all with ease and grace! 🌟

## License

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

## Copyright

Copyright © 2023 Daniel S. (GitHub: periaate)
All rights reserved.

## Authors
- Daniel S. (GitHub: periaate)
  
## Acknowledgements

A big shoutout to the amazing Go community and the developers who have contributed to the Go ecosystem. Your hard work and dedication make tools like the Blazing Fast Directory Lister possible. Thank you! 🚀