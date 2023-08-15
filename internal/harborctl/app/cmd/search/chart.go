/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package search

import (
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/client"
	"github.com/spf13/cobra"
	"os"
)

func ChartCmd(options *root.GlobalOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "chart",
		Short: "Search charts or project by fuzzy query , support many keywords separated by ',' (e.g. nginx,busy)",
		Long: `Usage:  harborctl search  chart pattern ;
For example:
# you can find somethings with one keyword!
  harborctl search chart nginx  OR
  harborctl search chart ng
# you can find somethings many keywords!
  harborctl search chart nginx,busybox  OR
  harborctl search chart ng,busy
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("has not enough analysis!!")
				os.Exit(1)
			}

			if err := client.NewChart(options).SearchChart(args[0]); err != nil {
				panic(err)
			}
		},
	}
	return command
}
