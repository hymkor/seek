seek - color-grep 
=================

[![GoDev](https://pkg.go.dev/badge/github.com/hymkor/seek)](https://pkg.go.dev/github.com/hymkor/seek)

```
Usage: seek [flags...] REGEXP Files...
  -A int
        print N lines after matching lines
  -B int
        print N lines before matching lines
  -html
        output html
  -i    ignore case
  -m value
        multi regular expression
  -no-color
        no color
  -r    recursive
```

On Windows

* If the line is valid as UTF8, seek.exe consider the line encoded by UTF8.
* If the line is invalid as UTF8, seek.exe conider the line encoded by ANSI(the encoding of the current code page.)

Install
-------

Download the binary package from [Releases](https://github.com/hymkor/seek/releases) and extract the executable.

### for scoop-installer

```
scoop install https://raw.githubusercontent.com/hymkor/seek/master/seek.json
```

or

```
scoop bucket add hymkor https://github.com/hymkor/scoop-bucket
scoop install seek
```
