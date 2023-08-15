/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package search

import (
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/client"
	"io"

	"github.com/spf13/cobra"
)

type projectOptions struct {
	global *root.GlobalOptions
	isAll  bool
}

// projectCmd represents the project command

func ProjectCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &projectOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "project",
		Short: "Search project by fuzzy query , support many keywords separated by ','",
		Long:  `search project very nice!`,
		RunE:  action.CommandAction(opts.run),
	}
	command.Flags().BoolVarP(&opts.isAll, "all", "a", false, "list projects")
	return command
}
func (c *projectOptions) run(args []string, stdout io.Writer) error {
	//if isAll is true then project list
	if c.isAll {
		return client.NewProject(c.global).SearchProjectsList()
	}

	//if isAll is false and args[0] has args then search project by name
	if len(args) != 0 {
		return client.NewProject(c.global).SearchProjects(args[0])
	}
	return fmt.Errorf("must provided one project name \n")
}
