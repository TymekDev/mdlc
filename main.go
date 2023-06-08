package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var version string

func main() {
	log.SetFlags(0)

	cmd := &cobra.Command{
		Use:     "mdlc [flags] file [...]",
		Short:   "mdlc - Markdown Link Checker",
		Version: version,
		// CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(aggregate(args))
		},
	}

	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
