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

type executionOptions struct {
	global *root.GlobalOptions
	reID   string
}

func ExecutionCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &executionOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "execution",
		Short: "Search execution or search by replicationID",
		Long: `For Example:
 harborctl search execution -i ID`,
		RunE: action.CommandAction(opts.run),
	}
	command.Flags().StringVarP(&opts.reID, "id", "i", "", "search by executionID")
	return command
}

func (c *executionOptions) run(args []string, stdout io.Writer) error {
	if c.reID != "" {
		client.NewReplication(c.global).SearchExecutionByID(c.reID)
	}
	return nil
}
