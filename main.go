package main

import (
	"log"

	"github.com/spf13/cobra"
)

var version string

const description = `mdlc - Markdown Link Checker

Description:
  mdlc scans markdown files for links and checks their status using a HTTP HEAD
  request. This includes checking both, status code and any redirects.

Notes:
  mdlc does not verify whether fragment URLs (starting with '#') are correct.`

func main() {
	log.SetFlags(0)

	cmd := &cobra.Command{
		Use:     "mdlc [flags] file [...]",
		Short:   description,
		Version: version,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := cmd.Flags().GetString("format")
			if err != nil {
				return err
			}
			flat, err := cmd.Flags().GetBool("flat")
			if err != nil {
				return err
			}
			return output(aggregate(args), format, flat)
		},
	}

	cmd.Flags().String("format", "columns", "output `format`: columns, json, or tsv")
	cmd.Flags().Bool("flat", false, "flatten JSON output to a single array")
	_ = cmd.Execute()
}
