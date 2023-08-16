package delete

import (
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

type deleteOptions struct {
	global *root.GlobalOptions
	id     int64
}

func DeleteCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &deleteOptions{
		global: options,
	}
	command := &cobra.Command{
		Use:   "delete",
		Short: "delete registry/replication/project",
		Long: `For example:
  harborctl delete replication -i replicationID 
  harborctl delete registry  -i registryID
  harborctl delete project NAME 
`,
		RunE: action.CommandAction(opts.run),
	}
	command.Flags().Int64VarP(&opts.id, "id", "i", 0, "when create replication must set this registryID")
	return command
}

func (d *deleteOptions) run(args []string, stdout io.Writer) error {

	if len(args) == 0 {
		return errors.New("no args must be registry or replication")
	}
	switch args[0] {
	case "registry":
		return client.NewRegistry(d.global).DeleteRegistry(d.id)
	case "replication":
		return client.NewReplication(d.global).DeleteReplication(d.id)
	case "project":
		if len(args) == 2 {
			return client.NewProject(d.global).DeleteProject(args[1])
		}
		return fmt.Errorf("project name must be provided")
	default:
		return errors.New("command execute failed ")
	}
}
