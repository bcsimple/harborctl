package main

import (
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/config"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/create"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/download"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/root"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/scan"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/search"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/start"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/update"
	"github.com/bcsimple/harborctl/internal/harborctl/app/cmd/version"
	"github.com/spf13/cobra"
	"os"
)

func createApp() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	var g = root.GlobalOptions{}
	rootCmd := &cobra.Command{
		Use:   "harborctl",
		Short: "A simple command tool for harbor",
		Long: `You can fast find something use this tool! help you do some tasks fastly!
create by zhangshun!`,
	}
	rootCmd.PersistentFlags().StringVarP(&g.Context, "context", "c", "", "set context for which harbor")
	rootCmd.PersistentFlags().StringVarP(&g.FormatStyle, "format", "f", "table", "set print format options table or kube)")
	rootCmd.AddCommand(
		create.CreateCmd(&g),
		search.SearchCmd(&g),
		download.DownloadCmd(&g),
		start.StartCmd(&g),
		update.UpdateCmd(&g),
		config.ConfigCmd(&g),
		scan.ScanCmd(&g),
		version.VersionCmd(),
		//copy.CopyCmd(),
	)
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := createApp().Execute(); err != nil {
		os.Exit(1)
	}
}

func main() {
	Execute()
}
