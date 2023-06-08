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

	cmd.Flags().String("format", "json", "output `format`: json or tsv")
	cmd.Flags().Bool("flat", false, "flatten JSON output to a single array")
	_ = cmd.Execute()
}
