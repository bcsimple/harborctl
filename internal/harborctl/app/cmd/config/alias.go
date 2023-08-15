package config

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/config"
	"github.com/spf13/cobra"
	"io"
)

type aliasOptions struct {
	global *root.GlobalOptions
	alias  string
}

func AliasCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &aliasOptions{}

	command := &cobra.Command{
		Use:   "alias",
		Short: "set alias for harbor context",
		Long: `For example:
 harborctl config alias name -a ALIAS`,
		RunE: action.CommandAction(opts.run),
		Args: cobra.ExactArgs(1),
	}
	command.Flags().StringVarP(&opts.alias, "alias", "a", "", "set alias")
	_ = command.MarkFlagRequired("alias")
	return command
}

func (a *aliasOptions) run(args []string, stdout io.Writer) error {
	if err := config.NewConnectConfiguration().SetConnectInfoAlias(args[0], a.alias); err != nil {
		return err
	}
	return nil
}
