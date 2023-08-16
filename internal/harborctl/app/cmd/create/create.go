/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

type createOptions struct {
	global              *root.GlobalOptions
	registryID          string
	registryName        string
	registryReplaceName string
	//create pull module policy
	isReplicationPullPolicy bool
}

// createCommand represents the create command
func CreateCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &createOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "create",
		Short: "Create registry/replication but create replication base on registry ",
		Long: `For example:
  harborctl create replication -i registryID 
  harborctl create registry 
  harborctl create project NAME 
`,
		RunE: action.CommandAction(opts.run),
	}
	command.Flags().StringVarP(&opts.registryID, "id", "i", "", "when create replication must set this registryID")
	command.Flags().StringVarP(&opts.registryName, "name", "n", "", "when create registry must set this name(or alias) from harbor config")
	command.Flags().StringVarP(&opts.registryReplaceName, "replace-name", "", "", "when create registry indicate a name for replace origin name ")
	command.Flags().BoolVarP(&opts.isReplicationPullPolicy, "pull", "p", false, "when create replication default pushed module,pulled module if true")
	return command
}

func (c *createOptions) run(args []string, stdout io.Writer) error {

	if len(args) == 0 {
		return errors.New("no args must be registry or replication")
	}
	switch args[0] {
	case "registry":
		if c.registryName == "" {
			return client.NewRegistry(c.global).CreateRegistry()
		}
		return client.NewRegistry(c.global).CreateRegistryByConfigInfo(c.registryName, c.registryReplaceName)
	case "replication":
		return client.NewReplication(c.global).CreateReplication(c.registryID, c.isReplicationPullPolicy)
	case "project":
		if len(args) == 2 {
			return client.NewProject(c.global).CreateProject(args[1])
		}
		return fmt.Errorf("project name must be provided")
	default:
		return errors.New("command execute failed ")
	}
}
