/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package search

import (
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/client"

	"github.com/spf13/cobra"
)

func ImageCmd(options *root.GlobalOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "image",
		Short: "Search images or project by fuzzy query , support many keywords separated by ',' (e.g. nginx,busy)",
		Long: `Usage:  harborctl search  pattern ;  

For example:
# you can find somethings with one keyword!
  harborctl search image nginx  OR
  harborctl search image ng
# you can find somethings many keywords!
  harborctl search image nginx,busybox  OR
  harborctl search image ng,busy
`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Println("has not enough analysis!!")
			}
			if err := client.NewImage(options).SearchAll(args[0]); err != nil {
				panic(err)
			}
		},
	}
	return command
}
