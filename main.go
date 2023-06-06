package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var (
	version  string
	fVersion bool
)

// TODO: add cobra for completions
func main() {
	log.SetFlags(0)

	if version != "" {
		flag.BoolVar(&fVersion, "version", false, "version for mdlsc")
	}
	flag.Parse()

	if fVersion {
		log.Println("mdlsc version", version)
		os.Exit(0)
	}

	all, err := parseFilesForLinks(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	// Check and count URLs
	unique := Links{}
	m := map[string]*Link{}
	for _, link := range all {
		if _, ok := m[link.URL]; ok {
			m[link.URL].Count++
			continue
		}

		unique = append(unique, link)
		m[link.URL] = link
		resp, err := http.Head(link.URL)
		if err != nil {
			link.Err = err
			continue
		}
		link.StatusCode = resp.StatusCode
		if url := resp.Request.URL.String(); url != link.URL {
			link.Err = fmt.Errorf("indirect URL to: %s", url)
		}
	}

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].Less(unique[j])
	})

	log.Println(unique)
}

func parseFilesForLinks(filenames []string) (Links, error) {
	result := Links{}
	for _, filename := range filenames {
		links, err := parseFileForLinks(filename)
		if err != nil {
			return nil, err
		}
		result = append(result, links...)
	}
	return result, nil
}

func parseFileForLinks(filename string) (Links, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	links := Links{}
	document := goldmark.DefaultParser().Parse(text.NewReader(b))
	if err := ast.Walk(document, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		// PERF: this could skip some nodes
		if link, ok := n.(*ast.Link); ok {
			links = append(links, &Link{FileName: filename, URL: string(link.Destination), Count: 1})
		}
		return ast.WalkContinue, nil
	}); err != nil {
		return nil, err
	}
	return links, nil
}
