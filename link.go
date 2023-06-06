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
	FileName   string
	URL        string
	Count      int
	StatusCode int
	Err        error
}

func (l *Link) String() string {
	parts := []string{
		l.FileName,
		l.URL,
		strconv.Itoa(l.Count),
		strconv.Itoa(l.StatusCode),
	}
	if l.Err != nil {
		parts = append(parts, l.Err.Error())
	}
	return strings.Join(parts, "\t")
}

func (l *Link) Less(other *Link) bool {
	return l.FileName < other.FileName || (l.FileName == other.FileName && l.URL < other.URL)
}
