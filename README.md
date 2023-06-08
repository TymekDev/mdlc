# mdlc

mdlc scans markdown files for links and checks their status using a HTTP HEAD request.
This includes checking both, status code and any redirects.

Note: mdlc does not verify whether fragment URLs (starting with '#') are correct.

## Installation
Run `make` to compile `mdlc`.
Run `make install` to install `mdlc` and completions to `/usr/local/`.
Clean up with `make clean` and `make uninstall`, respectively.

To override `/usr/local/` PREFIX variable use `make -e PREFIX=/foo/bar/baz/`.

## Usage
```
mdlc - Markdown Link Checker

Description:
  mdlc scans markdown files for links and checks their status using a HTTP HEAD
  request. This includes checking both, status code and any redirects.

Notes:
  mdlc does not verify whether fragment URLs (starting with '#') are correct.

Usage:
  mdlc [flags] file [...]

Flags:
      --flat            flatten JSON output to a single array
      --format format   output format: tsv or json (default "tsv")
  -h, --help            help for mdlc
  -v, --version         version for mdlc
```

## Name
\[M\]ark\[d\]own \[L\]ink \[C\]hecker.
