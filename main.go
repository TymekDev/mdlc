package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var (
	version  string
	fVersion bool
)

func main() {
	log.SetFlags(0)

	if version != "" {
		flag.BoolVar(&fVersion, "version", false, "version for mdlc")
	}
	flag.Parse()

	if fVersion {
		fmt.Printf("mdlc version %s\n", version)
		return
	}

	all, err := parseFilesForLinks(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	// Check and count URLs
	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	unique := Links{}
	m := map[string]int{}
	for _, link := range all {
		if idx, ok := m[link.URL]; ok {
			unique[idx].Count++
			continue
		}

		m[link.URL] = len(unique)
		unique = append(unique, link)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			link := unique[i]
			resp, err := http.Head(link.URL)
			defer mu.Unlock()
			mu.Lock()
			if err != nil {
				link.Err = err
				return
			}
			link.StatusCode = resp.StatusCode
			if url := resp.Request.URL.String(); url != link.URL {
				link.Err = fmt.Errorf("indirect URL to: %s", url)
			}
		}(len(unique) - 1)
	}
	wg.Wait()

	sort.Slice(unique, func(i, j int) bool {
		return unique[i].Less(unique[j])
	})

	fmt.Println(unique)
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
	if err := ast.Walk(document, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}
		link, ok := n.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}
		links = append(links, &Link{FileName: filename, URL: string(link.Destination), Count: 1})
		return ast.WalkSkipChildren, nil
	}); err != nil {
		return nil, err
	}
	return links, nil
}
