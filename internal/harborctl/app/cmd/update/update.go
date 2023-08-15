/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package update

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/client"
	"github.com/spf13/cobra"
	"io"
)

type updateOptions struct {
	global *root.GlobalOptions
	ruID   string
}

func UpdateCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &updateOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "update",
		Short: "Update  by replicationID",
		Long: `For Example:
 harborctl update replication -i ID`,
		RunE: action.CommandAction(opts.run),
	}
	command.Flags().StringVarP(&opts.ruID, "id", "i", "", "update by replicationID")
	command.MarkFlagRequired("id")
	return command
}

func (c *updateOptions) run(args []string, stdout io.Writer) error {
	if c.ruID != "" {
		client.NewReplication(c.global).ModifyReplication(c.ruID)
	}
	return nil
}
