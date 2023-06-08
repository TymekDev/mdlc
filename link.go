package main

import (
	"fmt"
	"net/http"
)

type Link struct {
	Filename    string `json:"filename"`
	Destination string `json:"destination"`
	Count       int    `json:"count"`
	StatusCode  int    `json:"status_code"`
	ErrMsg      string `json:"error,omitempty"`
}

// NOTE: check isn't a method to pre
func (l Link) check() (int, string) {
	resp, err := http.Head(l.Destination)
	if err != nil {
		return 0, err.Error()
	}
	if url := resp.Request.URL.String(); url != l.Destination {
		return resp.StatusCode, fmt.Sprintf("indirect URL to: %s", url)
	}
	return resp.StatusCode, ""
}
