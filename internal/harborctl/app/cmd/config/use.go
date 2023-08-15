package config

import (
	"fmt"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/config"
	"github.com/spf13/cobra"
	"io"
)

type useOptions struct {
	name string
}

func ContextCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &useOptions{}

	command := &cobra.Command{
		Use:   "use-context",
		Short: "use current context name for harbor",
		Long: `For example:
 harborctl config use-context name `,
		RunE: action.CommandAction(opts.run),
	}
	command.Flags().StringVarP(&opts.name, "name", "n", "", "set context name for connection")
	return command
}

func (u *useOptions) run(args []string, stdout io.Writer) error {
	name := u.name
	if name == "" {
		if len(args) == 0 {
			return fmt.Errorf("must be required name or alias")
		}
		name = args[0]
	}
	config.NewConnectConfiguration().SetConnectInfoContext(name)
	return nil
}
