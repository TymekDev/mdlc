# mdlc

mdlc scans markdown files for links and checks their status using a HTTP HEAD request.
This includes checking both, status code and any redirects.

Note: mdlc does not verify whether fragment URLs (starting with '#') are correct.

## Installation
Run `make` to compile mdlc.
Run `make install` to install mdlc and completions to `/usr/local/`.
Clean up with `make clean` and `make uninstall`, respectively.

To override `/usr/local/` PREFIX variable use `make -e PREFIX=/foo/bar/baz/`.

## Usage
```
mdlc - Markdown Link Checker

Description:
  mdlc scans markdown files for links and checks their status using a HTTP HEAD
  request. This includes checking both, status code and any redirects.

Notes:
  - mdlc does not verify whether fragment URLs (starting with '#') are correct.
  - '-' can be used as a file to read from stdin. Using it multiple times has
    no effect past the first occurence.

Usage:
  mdlc [flags] file [...]

Flags:
      --flat            flatten JSON output to a single array
      --format format   output format: columns, json, or tsv (default "columns")
  -h, --help            help for mdlc
  -v, --version         version for mdlc
```

## Examples
Given a following `README.md`:
```markdown
[byocss](https://sr.ht/~tymek/byocss) is a batteries excluded [Hugo](https://gohugo.io) theme.
It provides semantic HTML layout without any styling (see an example at [byocss.tymek.dev](https://byocss.tymek.dev)).
Only one CSS class and two IDs are used throughout byocss! IDs also double down as anchors for sections.
**If you know CSS and want to heavily customize a theme without building it from scratch, then byocss is for you!**
```
Running `mdlc README.md` yields the following output:
```
README.md	https://byocss.tymek.dev   	1	200	OK
README.md	https://gohugo.io          	1	200	OK
README.md	https://sr.ht/~tymek/byocss	1	200	Indirect URL to: https://sr.ht/~tymek/byocss/
```

mdlc handles link aliases too:
```markdown
Here is an [example link] and [another], twice: [another]. [This errors](incorrect#content). [This gets ignored](#content)

[example link]: https://example.com
[another]: https://example.com/another
```
Output:
```
-       #content                        1       0       Skip: fragment URL
-       https://example.com             1       200     OK
-       https://example.com/another     2       404     Not Found
-       incorrect#content               1       0       Head "incorrect#content": unsupported protocol scheme ""
```

Note: '-' come from the input being passed via stdin.

## Name
\[M\]ark\[d\]own \[L\]ink \[C\]hecker.
