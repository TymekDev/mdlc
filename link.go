package main

import (
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
	Filename   string
	URL        string
	Count      int
	StatusCode int
	Err        error
}

func (l *Link) String() string {
	parts := []string{
		l.Filename,
		l.URL,
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
	return l.Filename < other.Filename || (l.Filename == other.Filename && l.URL < other.URL)
}
