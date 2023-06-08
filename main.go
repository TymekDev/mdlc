package main

import (
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
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			return output(aggregate(args), f)
		},
	}

	cmd.Flags().String("format", "json", "output `format`: json or tsv")
	_ = cmd.Execute()
}
