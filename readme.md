seek - color-grep 
=================

```
Usage: seek [flags...] REGEXP Files...
  -html
        output html
  -i    ignore case
  -r    recursive
```

On Windows

* If the line is valid as UTF8, seek.exe consider the line encoded by UTF8.
* If the line is invalid as UTF8, seek.exe conider the line encoded by ANSI(the encoding of the current code page.)
