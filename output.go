package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

func output(m map[string]map[string]*Link, format string, flat bool) error {
	var colFilename, colDestination int
	switch format {
	case "json":
		if flat {
			return json.NewEncoder(os.Stdout).Encode(flatten(m))
		}
		return json.NewEncoder(os.Stdout).Encode(m)
	case "columns":
		for filename, links := range m {
			if n := len(filename); n > colFilename {
				colFilename = n
			}
			for destination := range links {
				if n := len(destination); n > colDestination {
					colDestination = n
				}
			}
		}
		fallthrough
	case "tsv":
		links := flatten(m)
		sort.Slice(links, func(i, j int) bool { // sort by filename and then by destination
			return links[i].Filename < links[j].Filename || (links[i].Filename == links[j].Filename && links[i].Destination < links[j].Destination)
		})
		for _, l := range links {
			fmt.Printf("%-*s\t%-*s\t%d\t%d\t%s\n", colFilename, l.Filename, colDestination, l.Destination, l.Count, l.StatusCode, l.ErrMsg)
		}
	default:
		return fmt.Errorf("unsupported format: '%s'", format)
	}
	return nil
}

func flatten(m map[string]map[string]*Link) []*Link {
	result := []*Link{}
	for _, v := range m {
		for _, l := range v {
			result = append(result, l)
		}
	}
	return result
}
