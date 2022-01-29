seek - color-grep 
=================

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
