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

// TODO:
// - przyjmuję listę plików
// - dla każdego pliku znajduję wszystkie linki
// - w zależności od trybu:
//     - count: tylko zlicza
//     - check: sprawdza status code
//     - vet: sprawdza czy linki prowadzą do redirectów
//     - domyślnie wszystkie: --check --count --vet
// - output:
//	  - posortowany (chyba, że --no-sort)
//	  - z podziałem na pliki (chyba, że --no-group ~do wymyślenia inna nazwa)
//    - format=json/pretty
// - epic:
//	  - --fix dla tych rzeczy, które --vet znajduje

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
		fmt.Printf("mdlsc version %s\n", version)
		return
	}

	all, err := parseFilesForLinks(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	// Check and count URLs
	unique := Links{}
	m := map[string]int{}
	for _, link := range all {
		if idx, ok := m[link.URL]; ok {
			unique[idx].Count++
			continue
		}

		m[link.URL] = len(unique)
		unique = append(unique, link)

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
