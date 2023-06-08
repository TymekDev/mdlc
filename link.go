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
	Err         error
}

func check(l *Link) (int, error) {
	resp, err := http.Head(l.Destination)
	if err != nil {
		return 0, err
	}
	if url := resp.Request.URL.String(); url != l.Destination {
		return resp.StatusCode, fmt.Errorf("indirect URL to: %s", url)
	}
	return resp.StatusCode, nil
}
