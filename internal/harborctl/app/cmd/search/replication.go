/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package search

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/client"
	"io"

	"github.com/spf13/cobra"
)

type replicationOptions struct {
	global *root.GlobalOptions
	rpID   string
}

func ReplicationCmd(options *root.GlobalOptions) *cobra.Command {

	opts := &replicationOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "replication",
		Short: "Search replication or search by replicationID",
		Long: `For Example:
 harborctl search replication pattern Or
 harborctl search replication -i ID`,
		RunE: action.CommandAction(opts.run),
	}
	command.Flags().StringVarP(&opts.rpID, "id", "i", "", "search by replicationID")
	return command

}

func (c *replicationOptions) run(args []string, stdout io.Writer) error {
	if len(args) != 0 {
		client.NewReplication(c.global).SearchReplication(args[0])
	} else if c.rpID != "" {
		client.NewReplication(c.global).SearchReplicationByID(c.rpID, true)
	}
	return nil
}
