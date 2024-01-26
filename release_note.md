v0.2.0
=======
Jan 29, 2022

- If output is not a terminal, disable colors automatically
- `-m` option: Multi regular expressions
- Implement `-no-color` option and set using color default
- Implement `-A` and `-B` option
- Do not warn that an argument is directory
- Fix: When `seek hoge *.frm` & \*.frm do not exist, seek.exe waited stdin
- Use mattn/go-zglob instead of original wildcard expanding package
- Move wildcard expansion all to starstar/
- starstar: add callback version
- Support UTF16 by go-texts/mbcs
- Add option -html: output with html
- Support Linux
- Replace `\<` and `\>` to `\b`
- If STDOUT is not terminal, cut ANSI-Escape Sequence
- Support \*\*/\*.go : find files recursively
- exit as ERRORLEVEL=1 if not found.
- argf: wildcard matches with directory

v0.1.0
=======
Nov 28, 2017

The first release
