package config

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/config"
	"github.com/spf13/cobra"
	"io"
)

type viewOptions struct {
}

func ViewCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &viewOptions{}

	command := &cobra.Command{
		Use:   "view",
		Short: "print all harbor config",
		Long: `For example:
 harborctl config view`,
		RunE: action.CommandAction(opts.run),
	}
	return command
}

func (u *viewOptions) run(args []string, stdout io.Writer) error {
	config.NewConnectConfiguration().View()
	return nil
}

type printOptions struct {
	completeInfo bool
}

func PrintCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &printOptions{}

	command := &cobra.Command{
		Use:   "current-context",
		Short: "print current-context harbor config",
		Long: `For example:
 harborctl config current-context`,
		RunE: action.CommandAction(opts.run),
	}
	command.Flags().BoolVarP(&opts.completeInfo, "all", "a", false, "print complete info")
	return command
}

func (u *printOptions) run(args []string, stdout io.Writer) error {
	config.NewConnectConfiguration().PrintCurrentContext(u.completeInfo)
	return nil
}

type listOptions struct {
	global    *root.GlobalOptions
	onlyName  bool
	isDecrypt bool
}

func ListCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &listOptions{
		global: options,
	}

	command := &cobra.Command{
		Use:   "list",
		Short: "list all harbor config with format",
		Long: `For example:
 harborctl config list`,
		RunE: action.CommandAction(opts.run),
	}
	command.Flags().BoolVarP(&opts.onlyName, "onlyname", "o", false, "print name only")
	command.Flags().BoolVarP(&opts.isDecrypt, "isDecrypt", "d", false, "password decrypt")
	return command
}

func (u *listOptions) run(args []string, stdout io.Writer) error {
	config.NewConnectConfiguration().List(u.global.FormatStyle, u.onlyName, u.isDecrypt)
	return nil
}
