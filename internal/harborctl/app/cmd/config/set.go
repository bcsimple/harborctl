package config

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/pkg/action"
	"github.com/bcsimple/harborctl/pkg/config"
	"github.com/spf13/cobra"
	"io"
)

type setOptions struct {
	global *root.GlobalOptions
	*config.HarborConnectInfo
}

func SetCmd(options *root.GlobalOptions) *cobra.Command {
	opts := &setOptions{
		HarborConnectInfo: &config.HarborConnectInfo{},
	}
	command := &cobra.Command{
		Use:   "set-context",
		Short: "set current context name for harbor",
		Long: `For example:
 harborctl config set NAME -u USER -p PASSWORD -s HOST -a ALIAS`,
		RunE: action.CommandAction(opts.run),
		Args: cobra.ExactArgs(1),
	}
	command.Flags().StringVarP(&opts.HarborConnectInfo.Host, "server", "s", "", "set context name for connection and must set format: HOST:PORT")
	command.Flags().StringVarP(&opts.HarborConnectInfo.Name, "name", "n", "", "set context name for connection")
	command.Flags().StringVarP(&opts.HarborConnectInfo.User, "user", "u", "", "set context name for connection")
	command.Flags().StringVarP(&opts.HarborConnectInfo.Password, "password", "p", "", "set context name for connection")
	command.Flags().StringArrayVarP(&opts.HarborConnectInfo.Alias, "alias", "a", []string{}, "set context name for connection")
	command.Flags().MarkHidden("name")
	//_ = command.MarkFlagRequired("host")
	//_ = command.MarkFlagRequired("user")
	//_ = command.MarkFlagRequired("password")
	//_ = command.MarkFlagRequired("name")
	return command
}

func (u *setOptions) run(args []string, stdout io.Writer) error {
	if u.HarborConnectInfo.Name == "" {
		u.HarborConnectInfo.Name = args[0]
	}
	config.NewConnectConfiguration().SetConnectInfo(u.HarborConnectInfo)
	return nil
}
