package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Links []*Link

func (l Links) String() string {
	var sb strings.Builder
	for i, link := range l {
		sb.WriteString(link.String())
		if i < len(l)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

type Link struct {
	Filename    string
	Destination string
	Count       int
	StatusCode  int
	Err         error
}

func (l *Link) String() string {
	parts := []string{
		l.Filename,
		l.Destination,
		strconv.Itoa(l.Count),
		strconv.Itoa(l.StatusCode),
	}
	if l.Err != nil {
		parts = append(parts, l.Err.Error())
	} else {
		parts = append(parts, "OK")
	}
	return strings.Join(parts, "\t")
}

func (l *Link) Less(other *Link) bool {
	return l.Filename < other.Filename || (l.Filename == other.Filename && l.Destination < other.Destination)
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
