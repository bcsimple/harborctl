package config

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/config"
	"github.com/spf13/cobra"
	"io"
)

type deleteOptions struct {
}

func DeleteCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &deleteOptions{}

	command := &cobra.Command{
		Use:   "delete-context",
		Short: "use current context name for harbor",
		Long: `For example:
 harborctl config delete-context name `,
		RunE: action.CommandAction(opts.run),
		Args: cobra.ExactArgs(1),
	}
	return command
}

func (u *deleteOptions) run(args []string, stdout io.Writer) error {
	config.NewConnectConfiguration().DelConnectInfo(args[0])
	return nil
}
