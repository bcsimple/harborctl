/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package start

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/client"
	"github.com/spf13/cobra"
	"io"
)

type startOptions struct {
	global     *root.GlobalOptions
	policyID   string
	policyPath string
}

func StartCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &startOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "start",
		Short: "Start replication running by replicationID",
		Long: `For example: 
  harborctl start -i number
`,
		RunE: action.CommandAction(opts.run),
	}

	command.Flags().StringVarP(&opts.policyID, "id", "i", "", "set a replicationID ")
	command.Flags().StringVarP(&opts.policyPath, "path", "p", "", "set a config path ")
	command.MarkFlagRequired("id")
	return command
}

func (c *startOptions) run(args []string, stdout io.Writer) error {
	if c.policyID != "" {
		if c.policyPath == "" {
			client.NewReplication(c.global).StartExecution(c.policyID)
		} else {
			client.NewReplication(c.global).StartExecutionFromConfig(c.policyID, c.policyPath)
		}
	}
	return nil
}
