package main

import (
	"fmt"
	"net/http"
)

type Link struct {
	Filename    string
	Destination string
	Count       int
	StatusCode  int
	ErrMsg      string
}

// NOTE: check isn't a method to pre
func (l Link) check() (int, string) {
	resp, err := http.Head(l.Destination)
	if err != nil {
		return 0, ""
	}
	if url := resp.Request.URL.String(); url != l.Destination {
		return resp.StatusCode, fmt.Sprintf("indirect URL to: %s", url)
	}
	return resp.StatusCode, ""
}
