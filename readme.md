seek - color-grep for ansi/utf8 (Windows)
=========================================

```
Usage: seek.exe [flags...] REGEXP Files...
  -i    ignore case
  -r    recursive
```

* If the line is valid as UTF8, seek.exe consider the line encoded by UTF8.
* If the line is invalid as UTF8, seek.exe conider the line encoded by ANSI(the encoding of the current code page.)
