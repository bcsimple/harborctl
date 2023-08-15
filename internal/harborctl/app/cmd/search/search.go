/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package search

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/spf13/cobra"
)

type searchOptions struct {
	global *root.GlobalOptions
}

func SearchCmd(options *root.GlobalOptions) *cobra.Command {

	command := &cobra.Command{
		Use:   "search",
		Short: "Search some resources  ",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	command.AddCommand(RegistryCmd(options))
	command.AddCommand(ReplicationCmd(options))
	command.AddCommand(TaskCmd(options))
	command.AddCommand(ExecutionCmd(options))
	command.AddCommand(ChartCmd(options))
	command.AddCommand(ProjectCmd(options))
	command.AddCommand(ImageCmd(options))
	return command
}
