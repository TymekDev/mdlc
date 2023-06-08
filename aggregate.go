package main

import (
	"log"
	"os"
	"sync"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

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
			mu.Lock()
			if _, ok := m[l.Filename]; !ok { // initalize map
				m[l.Filename] = map[string]*Link{}
			}
			if _, ok := m[l.Filename][l.Destination]; !ok { // increment existing link's count
				m[l.Filename][l.Destination] = l
			}
			m[l.Filename][l.Destination].Count++
			mu.Unlock()

			// Insert a new link
			sc, err := check(l) // PERF: we could cache responses in case one link appears in multiple files
			mu.Lock()
			l.StatusCode = sc
			l.Err = err
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
			ch <- &Link{Filename: filename, Destination: string(link.Destination)}
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
}
