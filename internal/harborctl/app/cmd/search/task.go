/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package search

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/client"
	"github.com/spf13/cobra"
	"io"
)

type taskOptions struct {
	global *root.GlobalOptions
	eID    string
	eSize  string
}

// taskCmd represents the tasks command

func TaskCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &taskOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "tasks",
		Short: "Search tasks or search by executionID",
		Long: `For Example:
 harborctl search tasks -i ID`,
		RunE: action.CommandAction(opts.run),
	}
	command.Flags().StringVarP(&opts.eID, "id", "i", "", "search by executionID")
	command.Flags().StringVarP(&opts.eSize, "number", "n", "", "setup tasks number default 10")
	command.MarkFlagRequired("id")
	return command
}

func (c *taskOptions) run(args []string, stdout io.Writer) error {
	if c.eID != "" {
		if c.eSize == "" {
			c.eSize = "10"
		}
		client.NewReplication(c.global).SearchTasksByID(c.eID, c.eSize)
	}
	return nil
}
