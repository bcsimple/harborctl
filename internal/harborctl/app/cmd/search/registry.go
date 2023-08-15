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

type registryOptions struct {
	global *root.GlobalOptions
	rID    string
}

func RegistryCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &registryOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "registry",
		Short: "Search registry or search by registryID",
		Long: `For Example:
 harborctl search registry pattern Or
 harborctl search registry -i ID`,
		RunE: action.CommandAction(opts.run),
	}

	command.Flags().StringVarP(&opts.rID, "id", "i", "", "search by registryID")
	return command
}

func (c *registryOptions) run(args []string, stdout io.Writer) error {
	if len(args) != 0 {
		client.NewRegistry(c.global).SearchRegistry(args[0])
	} else if c.rID != "" {
		client.NewRegistry(c.global).SearchRegistryByID(c.rID, true)
	}
	return nil
}
