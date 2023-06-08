package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

	fmt.Println(aggregate(flag.Args()))
}

func aggregate(filenames []string) map[string]map[string]*Link {
	ch := make(chan *Link)
	go collect(ch, filenames)

	m := map[string]map[string]*Link{}
	mu := sync.Mutex{} // PERF: we could have a mutex per inner map lock instead of a global one
	wg := sync.WaitGroup{}
	for link := range ch {
		wg.Add(1)
		go func(l *Link) {
			defer wg.Done()
			if _, ok := m[l.Filename]; !ok {
				mu.Lock()
				m[l.Filename] = map[string]*Link{}
				mu.Unlock()
			} else if _, ok := m[l.Filename][l.URL]; ok {
				mu.Lock()
				m[l.Filename][l.URL].Count++
				mu.Unlock()
				return
			}

			sc, err := check(l)
			l.Count = 1
			l.StatusCode = sc
			l.Err = err
			mu.Lock()
			m[l.Filename][l.URL] = l
			mu.Unlock()
		}(link)
	}
	wg.Wait()

	return m
}

func collect(ch chan *Link, filenames []string) {
	wg := sync.WaitGroup{}
	for _, filename := range filenames {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			readAndTraverse(ch, filename)
		}(filename)
	}
	wg.Wait()
	close(ch) // Close once all files have been read
}

func readAndTraverse(ch chan *Link, filename string) {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return
	}

	document := goldmark.DefaultParser().Parse(text.NewReader(b))
	_ = ast.Walk(document, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if !enter {
			return ast.WalkContinue, nil
		}
		link, ok := n.(*ast.Link)
		if ok {
			ch <- &Link{Filename: filename, URL: string(link.Destination)}
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
}
