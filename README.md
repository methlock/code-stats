*This simple CLI tool will scan path/paths for files and analyze 
them for some code statistics, like number of files, code rows, comment 
rows, extensions, etc*

Instalation
-
```bash
git clone https://github.com/methlock/code-stats.git
sudo cp code-stats/bin/code-stats /usr/local/bin
rm -r code-stats   
```

Type `code-stats help` for description.

Output:
```
Usage:
  code-stats [command]

Available Commands:
  help        Help about any command
  scan        Scan path or paths for files and process them

Flags:
  -h, --help   help for code-stats

Use "code-stats [command] --help" for more information about a command.
```

Usage
-
Type `$ code-stats scan -help` for usage and flags.

Output:
```
Scan path or paths for files and process them

Usage:
  code-stats scan [PATH or PATHS] [flags]

Flags:
  -x, --allowed-exts string   Specify the extensions to process. Pass them without dot, like '-x "go py"'
  -e, --exclude-path string   Paths to exclude. Usage: '-e "path1 path2"'
  -h, --help                  help for scan
  -s, --hidden                Search and process also hidden folders and files
  -p, --print-exts            Prints also information by extensions
  -f, --print-files           Prints stats for each file
  -n, --process-no-ext        Process also files with no extensions, like binaries
```

Example
-
Example usage on this repository.

`code-stats path/to/this/repo -pf`

Prints overall info as well as extensions and individual files stats.

Output:
```
Overall stats:
 Number of files: 7
 Total lines: 450
 Code lines: 416
 Comment lines: 34

Stats by extensions
Extension  Files  Total lines  Code lines  Comment lines
md         1      81           81          0
go         6      369          335         34

Stats by files:
Name       Total lines  Code lines  Comment lines  Path
README.md  81           81          0              .../code-stats/README.md
root.go    26           23          3              .../code-stats/cmd/root.go
scan.go    214          193         21             .../code-stats/cmd/scan.go
config.go  6            5           1              .../code-stats/config/config.go
main.go    29           29          0              .../code-stats/main.go
models.go  54           50          4              .../code-stats/models/models.go
utils.go   40           35          5              .../code-stats/utils/utils.go
```