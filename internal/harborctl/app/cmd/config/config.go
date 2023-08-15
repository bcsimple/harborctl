package config

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/spf13/cobra"
)

func ConfigCmd(options *root.GlobalOptions) *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "config harbor context when we use harborctl ",
		Long: `For example:
 harborctl config set NAME -u USER -p PASSWORD -h HOST -a ALIAS 
 harborctl config use-context name
 harborctl config alias name -a ALIAS
`,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	command.AddCommand(SetCmd(options))
	command.AddCommand(ContextCmd(options))
	command.AddCommand(AliasCmd(options))
	command.AddCommand(ViewCmd(options))
	command.AddCommand(PrintCmd(options))
	command.AddCommand(DeleteCmd(options))
	command.AddCommand(ListCmd(options))
	return command
}
